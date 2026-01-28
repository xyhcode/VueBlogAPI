/*
 * @Description: Essay handler for HTTP requests
 * @Author: Qwenjie
 * @Date: 2026-01-27
 * @LastEditTime: 2026-01-27
 * @LastEditors: Qwenjie
 */
package essay_handler

import (
	"net/http"
	"strconv"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/response"
	"github.com/anzhiyu-c/anheyu-app/pkg/service/essay"
	"github.com/gin-gonic/gin"
)

// Handler encapsulates essay related controller methods
type Handler struct {
	essaySvc essay.Service
}

// NewHandler is the constructor for Handler
func NewHandler(essaySvc essay.Service) *Handler {
	return &Handler{
		essaySvc: essaySvc,
	}
}

// GetAllEssays handles getting all essay records with pagination
// @Summary      获取所有随笔记录（分页）
// @Description  获取所有随笔记录，无需认证，支持分页
// @Tags         随笔管理
// @Produce      json
// @Param        page     query  int  false  "页码"  default(1)
// @Param        pageSize query  int  false  "每页数量"  default(10)
// @Success      200  {object}  response.Response{data=object{list=[]object{id=integer,content=string,date=string,images=[]object{url=string,alt=string},link=string},total=integer,pageNum=integer,pageSize=integer}}  "获取成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /public/essay [get]
func (h *Handler) GetAllEssays(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务获取分页数据
	pageResult, err := h.essaySvc.GetEssaysByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取随笔记录失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":     pageResult.Items,
		"total":    pageResult.Total,
		"pageNum":  page,
		"pageSize": pageSize,
	}, "获取随笔记录成功")
}

// GetPublicEssays handles getting all public essay records (most recent first)
// @Summary      获取公开随笔记录
// @Description  获取所有公开随笔记录，按创建时间倒序排列，无需认证
// @Tags         随笔管理
// @Produce      json
// @Success      200  {object}  response.Response{data=object{list=[]object{id=integer,content=string,date=string,images=[]object{url=string,alt=string},link=string}}}  "获取成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /public/essays [get]
func (h *Handler) GetPublicEssays(c *gin.Context) {
	// 调用服务获取所有数据
	essays, err := h.essaySvc.GetAllEssays(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取随笔记录失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list": essays,
	}, "获取随笔记录成功")
}

// CreateEssay handles creating a new essay record
// @Summary      创建随笔记录
// @Description  创建新的随笔记录，需要认证
// @Tags         随笔管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  object{content=string,date=string,images=[]object{url=string,alt=string},link=string}  true  "随笔记录信息"
// @Success      200  {object}  response.Response{data=object{id=integer,content=string,date=string,images=[]object{url=string,alt=string},link=string}}  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "创建失败"
// @Router       /essay [post]
func (h *Handler) CreateEssay(c *gin.Context) {
	var req struct {
		Content string        `json:"content" binding:"required"`
		Date    string        `json:"date" binding:"required"`
		Images  []model.Image `json:"images,omitempty"` // JSON array
		Link    string        `json:"link,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	record, err := h.essaySvc.CreateEssay(c.Request.Context(), essay.CreateEssayParams{
		Content: req.Content,
		Date:    req.Date,
		Images:  req.Images, // Pass as JSON array
		Link:    req.Link,
	})

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建随笔记录失败: "+err.Error())
		return
	}

	response.Success(c, record, "创建随笔记录成功")
}

// UpdateEssay handles updating an existing essay record
// @Summary      更新随笔记录
// @Description  更新随笔记录，需要认证
// @Tags         随笔管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int     true  "记录ID"
// @Param        body  body  object{content=string,date=string,images=[]object{url=string,alt=string},link=string}  true  "随笔记录信息"
// @Success      200  {object}  response.Response{data=object{id=integer,content=string,date=string,images=[]object{url=string,alt=string},link=string}}  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误或ID非法"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "记录不存在"
// @Failure      500  {object}  response.Response  "更新失败"
// @Router       /essay/{id} [put]
func (h *Handler) UpdateEssay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
		return
	}

	var req struct {
		Content string        `json:"content" binding:"required"`
		Date    string        `json:"date" binding:"required"`
		Images  []model.Image `json:"images,omitempty"` // JSON array
		Link    string        `json:"link,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	record, err := h.essaySvc.UpdateEssay(c.Request.Context(), uint(id), essay.UpdateEssayParams{
		Content: req.Content,
		Date:    req.Date,
		Images:  req.Images, // Pass as JSON array
		Link:    req.Link,
	})

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新随笔记录失败: "+err.Error())
		return
	}

	response.Success(c, record, "更新随笔记录成功")
}

// DeleteEssay handles deleting an essay record
// @Summary      删除随笔记录
// @Description  删除随笔记录，需要认证
// @Tags         随笔管理
// @Security     BearerAuth
// @Param        id  path  int  true  "记录ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "ID非法"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "记录不存在"
// @Failure      500  {object}  response.Response  "删除失败"
// @Router       /essay/{id} [delete]
func (h *Handler) DeleteEssay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
		return
	}

	if err := h.essaySvc.DeleteEssay(c.Request.Context(), uint(id)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除随笔记录失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除随笔记录成功")
}
