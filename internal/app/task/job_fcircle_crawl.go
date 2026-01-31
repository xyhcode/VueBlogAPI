package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/anzhiyu-c/anheyu-app/ent"
	"github.com/anzhiyu-c/anheyu-app/pkg/crawler"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
	"github.com/anzhiyu-c/anheyu-app/pkg/service/fcircle"
	"github.com/redis/go-redis/v9"
)

// FCircleCrawlJob 爬取朋友圈文章的任务
type FCircleCrawlJob struct {
	logger     *slog.Logger
	db         *ent.Client
	linkRepo   repository.LinkRepository
	redis      *redis.Client
	fcircleSvc *fcircle.Service
}

// NewFCircleCrawlJob 创建新的朋友圈爬取任务
func NewFCircleCrawlJob(logger *slog.Logger, db *ent.Client, linkRepo repository.LinkRepository, redis *redis.Client) *FCircleCrawlJob {
	return &FCircleCrawlJob{
		logger:     logger,
		db:         db,
		linkRepo:   linkRepo,
		redis:      redis,
		fcircleSvc: fcircle.NewService(db, redis),
	}
}

// Run 执行爬取任务
func (j *FCircleCrawlJob) Run() {
	j.logger.Info("开始执行朋友圈爬取任务")
	start := time.Now()

	// 统计变量
	friendsNum := 0
	errorNum := 0

	// 1. 获取所有友链（不限制状态）
	j.logger.Info("开始获取所有友链")
	links, total, err := j.linkRepo.List(context.Background(), &model.ListLinksRequest{
		PaginationInput: model.PaginationInput{
			Page:     1,
			PageSize: 1000, // 一次性获取所有友链
		},
	})
	if err != nil {
		j.logger.Error("获取友链列表失败", slog.Any("error", err))
		errorNum++
		// 继续执行，即使获取友链失败
	} else {
		friendsNum = len(links)
		j.logger.Info("获取友链列表成功", slog.Int("count", friendsNum), slog.Int("total", total))

		// 2. 转换为爬虫需要的格式
		friends := make([]crawler.Friend, 0, friendsNum)
		for _, link := range links {
			friends = append(friends, crawler.Friend{
				Name:   link.Name,
				Link:   link.URL,
				Avatar: link.Logo,
				Descr:  link.Description,
			})
		}

		// 3. 初始化爬虫
		c, err := crawler.NewCrawler()
		if err != nil {
			j.logger.Error("初始化爬虫失败", slog.Any("error", err))
			errorNum++
			// 继续执行
		} else {
			// 4. 爬取所有友链的文章
			posts, crawlErrorNum, err := c.CrawlAllFriends(friends)
			if err != nil {
				j.logger.Error("爬取文章失败", slog.Any("error", err))
				errorNum += crawlErrorNum
				// 继续执行，即使部分爬取失败
			}

			j.logger.Info("爬取文章完成", slog.Int("count", len(posts)))

			// 5. 保存爬取的文章
			if len(posts) > 0 {
				err = j.fcircleSvc.SavePosts(posts)
				if err != nil {
					j.logger.Error("保存文章失败", slog.Any("error", err))
					errorNum++
					// 继续执行
				} else {
					j.logger.Info("保存文章成功")
				}
			}
		}
	}

	// 6. 清理过期文章
	err = j.fcircleSvc.CleanupExpiredPosts()
	if err != nil {
		j.logger.Error("清理过期文章失败", slog.Any("error", err))
		errorNum++
		// 继续执行
	} else {
		j.logger.Info("清理过期文章成功")
	}

	// 7. 更新统计信息
	if err := j.fcircleSvc.UpdateStatistic(friendsNum, errorNum); err != nil {
		j.logger.Error("更新统计信息失败", slog.Any("error", err))
	} else {
		j.logger.Info("更新统计信息成功", slog.Int("friends_num", friendsNum), slog.Int("error_num", errorNum))
	}

	j.logger.Info("朋友圈爬取任务执行完成", slog.Duration("duration", time.Since(start)))
	return
}
