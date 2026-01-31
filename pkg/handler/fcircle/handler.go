package fcircle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
	"github.com/anzhiyu-c/anheyu-app/pkg/response"
	"github.com/anzhiyu-c/anheyu-app/pkg/service/fcircle"
)

// Handler 处理朋友圈相关的请求
type Handler struct {
	fcircleSvc *fcircle.Service
	redis      *redis.Client
	linkRepo   repository.LinkRepository
}

// NewHandler 创建新的朋友圈处理器
func NewHandler(fcircleSvc *fcircle.Service, redis *redis.Client, linkRepo repository.LinkRepository) *Handler {
	return &Handler{
		fcircleSvc: fcircleSvc,
		redis:      redis,
		linkRepo:   linkRepo,
	}
}

// ListAll 获取完整统计信息与文章列表
func (h *Handler) ListAll(c *gin.Context) {
	// 解析参数
	startStr := c.DefaultQuery("start", "0")
	endStr := c.DefaultQuery("end", "-1")
	rule := c.DefaultQuery("rule", "created")

	// 验证参数
	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 {
		response.Fail(c, http.StatusBadRequest, "start error, please use valid integer")
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "end error, please use valid integer")
		return
	}

	if rule != "created" && rule != "updated" {
		response.Fail(c, http.StatusBadRequest, "rule error, please use 'created'/'updated'")
		return
	}

	// 获取统计信息
	statistic, err := h.fcircleSvc.GetStatistic()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	// 获取文章列表
	posts, err := h.fcircleSvc.GetPosts(start, end, rule)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	// 构建响应
	responseData := gin.H{
		"statistical_data": gin.H{
			"friends_num":       statistic.FriendsNum,
			"active_num":        statistic.ActiveNum,
			"error_num":         statistic.ErrorNum,
			"article_num":       statistic.ArticleNum,
			"last_updated_time": statistic.LastUpdatedTime.Format("2006-01-02 15:04:05"),
		},
		"article_data": make([]gin.H, 0, len(posts)),
	}

	// 添加文章数据
	for i, post := range posts {
		articleData := gin.H{
			"floor":       i + 1,
			"title":       post.Title,
			"created":     post.Created.Format("2006-01-02 15:04:05"),
			"updated":     post.Updated.Format("2006-01-02 15:04:05"),
			"link":        post.Link,
			"author":      post.Author,
			"avatar":      post.Avatar,
			"friend_link": post.FriendLink,
		}

		responseData["article_data"] = append(responseData["article_data"].([]gin.H), articleData)
	}

	response.Success(c, responseData, "获取朋友圈数据成功")
}

// GetRandomPost 获取随机文章
func (h *Handler) GetRandomPost(c *gin.Context) {
	// 解析参数
	numStr := c.DefaultQuery("num", "1")
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > 100 {
		response.Fail(c, http.StatusBadRequest, "num error, please use integer between 1 and 100")
		return
	}

	// 获取随机文章
	posts, err := h.fcircleSvc.GetRandomPosts(num)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取随机文章失败")
		return
	}

	if len(posts) == 0 {
		response.Fail(c, http.StatusNotFound, "not found")
		return
	}

	// 构建响应
	responseData := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		responseData = append(responseData, gin.H{
			"title":     post.Title,
			"created":   post.Created.Format("2006-01-02 15:04:05"),
			"updated":   post.Updated.Format("2006-01-02 15:04:05"),
			"link":      post.Link,
			"author":    post.Author,
			"avatar":    post.Avatar,
			"rule":      post.Rules,
			"createdAt": post.CrawledAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, responseData, "获取随机文章成功")
}

// GetFriendPosts 获取指定朋友文章列表
func (h *Handler) GetFriendPosts(c *gin.Context) {
	// 解析参数
	link := c.Query("link")
	numStr := c.DefaultQuery("num", "-1")
	rule := c.DefaultQuery("rule", "created")

	// 验证参数
	num, err := strconv.Atoi(numStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "num error, please use valid integer")
		return
	}

	if rule != "created" && rule != "updated" {
		response.Fail(c, http.StatusBadRequest, "rule error, please use 'created'/'updated'")
		return
	}

	// 获取文章列表
	posts, err := h.fcircleSvc.GetPostsByLink(link, num, rule)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	if len(posts) == 0 {
		response.Fail(c, http.StatusNotFound, "not found")
		return
	}

	// 构建响应
	responseData := gin.H{
		"statistical_data": gin.H{
			"name":        posts[0].Author,
			"link":        link,
			"avatar":      posts[0].Avatar,
			"article_num": len(posts),
		},
		"article_data": make([]gin.H, 0, len(posts)),
	}

	// 添加文章数据
	for i, post := range posts {
		articleData := gin.H{
			"floor":       i + 1,
			"title":       post.Title,
			"created":     post.Created.Format("2006-01-02 15:04:05"),
			"updated":     post.Updated.Format("2006-01-02 15:04:05"),
			"link":        post.Link,
			"author":      post.Author,
			"avatar":      post.Avatar,
			"friend_link": post.FriendLink,
		}

		responseData["article_data"] = append(responseData["article_data"].([]gin.H), articleData)
	}

	response.Success(c, responseData, "获取朋友文章列表成功")
}
