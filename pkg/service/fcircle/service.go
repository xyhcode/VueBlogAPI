package fcircle

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/anzhiyu-c/anheyu-app/ent"
	"github.com/anzhiyu-c/anheyu-app/ent/fcirclepost"
	"github.com/anzhiyu-c/anheyu-app/pkg/crawler"
	"github.com/redis/go-redis/v9"
)

// Service 处理朋友圈相关的业务逻辑
type Service struct {
	db    *ent.Client
	rng   *rand.Rand
	redis *redis.Client
}

// NewService 创建新的朋友圈服务
func NewService(db *ent.Client, redis *redis.Client) *Service {
	// 初始化随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return &Service{
		db:    db,
		rng:   rng,
		redis: redis,
	}
}

// SavePosts 保存爬取的文章
func (s *Service) SavePosts(posts []crawler.Post) error {
	ctx := context.Background()
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	// 收集所有文章链接
	links := make([]string, len(posts))
	for i, post := range posts {
		links[i] = post.Link
	}

	// 批量查询已存在的文章
	existingPosts, err := tx.FCirclePost.Query().
		Where(fcirclepost.LinkIn(links...)).
		All(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("批量查询文章失败: %w", err)
	}

	// 构建已存在文章的映射，便于快速查找
	existingMap := make(map[string]*ent.FCirclePost)
	for _, post := range existingPosts {
		existingMap[post.Link] = post
	}

	// 批量处理文章
	for _, post := range posts {
		// 标准化 friend_link，移除 trailing slash
		friendLink := post.FriendLink
		if len(friendLink) > 0 && friendLink[len(friendLink)-1] == '/' {
			friendLink = friendLink[:len(friendLink)-1]
		}

		// 解析时间字符串
		createdTime, _ := time.Parse("2006-01-02 15:04:05", post.Created)
		updatedTime, _ := time.Parse("2006-01-02 15:04:05", post.Updated)

		// 检查文章是否已存在
		if existing, ok := existingMap[post.Link]; ok {
			// 文章已存在，更新
			_, err = tx.FCirclePost.UpdateOne(existing).
				SetTitle(post.Title).
				SetCreated(createdTime).
				SetUpdated(updatedTime).
				SetAuthor(post.Author).
				SetAvatar(post.Avatar).
				SetFriendLink(friendLink).
				SetCrawledAt(time.Now()).
				Save(ctx)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("更新文章失败: %w", err)
			}
		} else {
			// 文章不存在，创建
			rules := ""
			if len(post.Rules) > 0 {
				rules = post.Rules[0]
			}

			_, err = tx.FCirclePost.Create().
				SetTitle(post.Title).
				SetLink(post.Link).
				SetCreated(createdTime).
				SetUpdated(updatedTime).
				SetAuthor(post.Author).
				SetAvatar(post.Avatar).
				SetFriendLink(friendLink).
				SetCrawledAt(time.Now()).
				SetRules(rules).
				Save(ctx)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("创建文章失败: %w", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 清除相关缓存
	// 1. 清除统计信息缓存
	s.redis.Del(ctx, "fcircle:statistic")
	// 2. 清除文章列表缓存（使用通配符）
	keys, _ := s.redis.Keys(ctx, "fcircle:posts:start:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}
	// 3. 清除随机文章缓存（使用通配符）
	keys, _ = s.redis.Keys(ctx, "fcircle:random:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}
	// 4. 清除友链文章缓存（使用通配符）
	keys, _ = s.redis.Keys(ctx, "fcircle:posts:link:*").Result()
	if len(keys) > 0 {
		s.redis.Del(ctx, keys...)
	}

	// 更新统计信息（默认值，实际值由爬取任务传入）
	if err := s.UpdateStatistic(0, 0); err != nil {
		return fmt.Errorf("更新统计信息失败: %w", err)
	}

	return nil
}

// CleanupExpiredPosts 清理过期的文章（超过60天）
func (s *Service) CleanupExpiredPosts() error {
	ctx := context.Background()
	cutoff := time.Now().AddDate(0, 0, -60)

	// 删除过期文章
	_, err := s.db.FCirclePost.Delete().
		Where(fcirclepost.CrawledAtLT(cutoff)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("清理过期文章失败: %w", err)
	}

	// 更新统计信息（默认值，实际值由爬取任务传入）
	if err := s.UpdateStatistic(0, 0); err != nil {
		return fmt.Errorf("更新统计信息失败: %w", err)
	}

	return nil
}

// UpdateStatistic 更新统计信息
func (s *Service) UpdateStatistic(friendsNum, errorNum int) error {
	ctx := context.Background()

	// 获取文章总数
	articleNum, err := s.db.FCirclePost.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("获取文章总数失败: %w", err)
	}

	// 获取活跃友链数量（当月有文章发布的友链）
	activeNum, err := s.getActiveFriendsCount()
	if err != nil {
		return fmt.Errorf("获取活跃友链数量失败: %w", err)
	}

	// 如果没有传入友链总数，从数据库获取
	if friendsNum == 0 {
		friendsNum, err = s.getDistinctFriendsCount()
		if err != nil {
			return fmt.Errorf("获取友链总数失败: %w", err)
		}
	}

	// 获取或创建统计信息
	statistic, err := s.db.FCircleStatistic.Query().First(ctx)
	if err != nil {
		// 创建新的统计信息
		_, err = s.db.FCircleStatistic.Create().
			SetFriendsNum(friendsNum).
			SetActiveNum(activeNum).
			SetErrorNum(errorNum).
			SetArticleNum(articleNum).
			SetLastUpdatedTime(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("创建统计信息失败: %w", err)
		}
	} else {
		// 更新现有统计信息
		_, err = s.db.FCircleStatistic.UpdateOne(statistic).
			SetFriendsNum(friendsNum).
			SetActiveNum(activeNum).
			SetErrorNum(errorNum).
			SetArticleNum(articleNum).
			SetLastUpdatedTime(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("更新统计信息失败: %w", err)
		}
	}

	// 清除统计信息缓存
	cacheKey := "fcircle:statistic"
	s.redis.Del(ctx, cacheKey)

	return nil
}

// getActiveFriendsCount 统计当月有文章发布的友链数量
func (s *Service) getActiveFriendsCount() (int, error) {
	ctx := context.Background()

	// 计算当月的起始时间
	now := time.Now()
	// 获取当月第一天
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 统计当月有文章发布的去重友链数量
	// 注意：这里使用Author字段作为友链标识，实际应该根据具体数据结构调整
	// 由于ent不直接支持Group by和Count distinct，我们需要先获取所有符合条件的文章，然后在内存中去重
	posts, err := s.db.FCirclePost.Query().
		Where(fcirclepost.CreatedGTE(firstDayOfMonth)).
		All(ctx)
	if err != nil {
		return 0, err
	}

	// 去重统计
	authorMap := make(map[string]bool)
	for _, post := range posts {
		if post.Author != "" {
			authorMap[post.Author] = true
		}
	}

	return len(authorMap), nil
}

// getDistinctFriendsCount 统计去重后的友链总数
func (s *Service) getDistinctFriendsCount() (int, error) {
	ctx := context.Background()

	// 统计去重后的友链数量
	// 注意：这里使用Author字段作为友链标识，实际应该根据具体数据结构调整
	// 由于ent不直接支持Group by和Count distinct，我们需要先获取所有文章，然后在内存中去重
	posts, err := s.db.FCirclePost.Query().
		All(ctx)
	if err != nil {
		return 0, err
	}

	// 去重统计
	authorMap := make(map[string]bool)
	for _, post := range posts {
		if post.Author != "" {
			authorMap[post.Author] = true
		}
	}

	return len(authorMap), nil
}

// GetStatistic 获取统计信息
func (s *Service) GetStatistic() (*ent.FCircleStatistic, error) {
	ctx := context.Background()

	// 尝试从 Redis 缓存获取
	cacheKey := "fcircle:statistic"
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var statistic ent.FCircleStatistic
		if err := json.Unmarshal([]byte(cachedData), &statistic); err == nil {
			return &statistic, nil
		}
	}

	// 缓存未命中，从数据库获取
	statistic, err := s.db.FCircleStatistic.Query().First(ctx)
	if err != nil {
		// 创建默认统计信息并更新
		if err := s.UpdateStatistic(0, 0); err != nil {
			return nil, fmt.Errorf("更新统计信息失败: %w", err)
		}
		// 重新查询统计信息
		statistic, err = s.db.FCircleStatistic.Query().First(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取统计信息失败: %w", err)
		}
	}

	// 更新缓存，设置 6 小时过期
	if data, err := json.Marshal(statistic); err == nil {
		s.redis.Set(ctx, cacheKey, data, 6*time.Hour)
	}

	return statistic, nil
}

// GetPosts 获取文章列表
func (s *Service) GetPosts(start, end int, rule string) ([]*ent.FCirclePost, error) {
	ctx := context.Background()

	// 生成缓存键
	cacheKey := fmt.Sprintf("fcircle:posts:start:%d:end:%d:rule:%s", start, end, rule)

	// 尝试从 Redis 缓存获取
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var posts []*ent.FCirclePost
		if err := json.Unmarshal([]byte(cachedData), &posts); err == nil {
			return posts, nil
		}
	}

	// 缓存未命中，从数据库获取
	query := s.db.FCirclePost.Query()

	// 排序
	if rule == "created" {
		query = query.Order(ent.Desc(fcirclepost.FieldCreated))
	} else {
		query = query.Order(ent.Desc(fcirclepost.FieldUpdated))
	}

	// 分页
	if start > 0 {
		query = query.Offset(start)
	}

	if end > start {
		query = query.Limit(end - start + 1)
	}

	posts, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取文章列表失败: %w", err)
	}

	// 更新缓存，设置 6 小时过期
	if data, err := json.Marshal(posts); err == nil {
		s.redis.Set(ctx, cacheKey, data, 6*time.Hour)
	}

	return posts, nil
}

// GetRandomPosts 获取随机文章
func (s *Service) GetRandomPosts(num int) ([]*ent.FCirclePost, error) {
	ctx := context.Background()

	// 从数据库获取所有文章
	allPosts, err := s.db.FCirclePost.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取文章失败: %w", err)
	}

	if len(allPosts) == 0 {
		return []*ent.FCirclePost{}, nil
	}

	// 随机打乱文章顺序
	for i := len(allPosts) - 1; i > 0; i-- {
		j := s.rng.Intn(i + 1)
		allPosts[i], allPosts[j] = allPosts[j], allPosts[i]
	}

	// 限制返回数量
	if num > len(allPosts) {
		num = len(allPosts)
	}

	return allPosts[:num], nil
}

// GetPostsByLink 根据链接获取文章
func (s *Service) GetPostsByLink(link string, num int, rule string) ([]*ent.FCirclePost, error) {
	ctx := context.Background()

	// 标准化 link，移除 trailing slash
	normalizedLink := link
	if len(normalizedLink) > 0 && normalizedLink[len(normalizedLink)-1] == '/' {
		normalizedLink = normalizedLink[:len(normalizedLink)-1]
	}

	// 生成缓存键
	cacheKey := fmt.Sprintf("fcircle:posts:link:%s:num:%d:rule:%s", normalizedLink, num, rule)

	// 尝试从 Redis 缓存获取
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var posts []*ent.FCirclePost
		if err := json.Unmarshal([]byte(cachedData), &posts); err == nil {
			return posts, nil
		}
	}

	// 缓存未命中，从数据库获取
	query := s.db.FCirclePost.Query()

	// 按作者过滤（假设 Author 字段标识朋友）
	if normalizedLink != "" {
		// 这里使用严格匹配，只有当 link 与友链用户网站链接完全匹配时才返回数据
		// 注意：在爬取时，我们也应该标准化 friend_link 字段，移除 trailing slash
		query = query.Where(fcirclepost.FriendLink(normalizedLink))
	} else {
		// 随机选择一个朋友
		// 1. 获取所有不同的友链用户网站链接
		allPosts, err := s.db.FCirclePost.Query().All(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取文章失败: %w", err)
		}

		// 收集所有不同的友链用户网站链接
		friendLinks := make(map[string]bool)
		for _, post := range allPosts {
			friendLinks[post.FriendLink] = true
		}

		// 随机选择一个友链用户网站链接
		friendLinkList := make([]string, 0, len(friendLinks))
		for friendLink := range friendLinks {
			friendLinkList = append(friendLinkList, friendLink)
		}

		if len(friendLinkList) > 0 {
			randomIndex := s.rng.Intn(len(friendLinkList))
			randomFriendLink := friendLinkList[randomIndex]
			query = query.Where(fcirclepost.FriendLink(randomFriendLink))
		}
	}

	// 排序
	if rule == "created" {
		query = query.Order(ent.Desc(fcirclepost.FieldCreated))
	} else {
		query = query.Order(ent.Desc(fcirclepost.FieldUpdated))
	}

	// 分页
	if num > 0 {
		query = query.Limit(num)
	}

	posts, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取文章列表失败: %w", err)
	}

	// 更新缓存，设置 6 小时过期
	if data, err := json.Marshal(posts); err == nil {
		s.redis.Set(ctx, cacheKey, data, 6*time.Hour)
	}

	return posts, nil
}
