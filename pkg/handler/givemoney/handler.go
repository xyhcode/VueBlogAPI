/*
 * @Description: GiveMoney handler for HTTP requests
 * @Author: Qwenjie
 * @Date: 2026-01-24
 * @LastEditTime: 2026-01-24
 * @LastEditors: Qwenjie
 */
package givemoney_handler

import (
	"net/http"
	"strconv"

	"github.com/anzhiyu-c/anheyu-app/pkg/response"
	"github.com/anzhiyu-c/anheyu-app/pkg/service/givemoney"
	"github.com/gin-gonic/gin"
)

// GiveMoneyHandler encapsulates give money related controller methods
type GiveMoneyHandler struct {
	giveMoneySvc givemoney.GiveMoneyService
}

// NewGiveMoneyHandler is the constructor for GiveMoneyHandler
func NewGiveMoneyHandler(giveMoneySvc givemoney.GiveMoneyService) *GiveMoneyHandler {
	return &GiveMoneyHandler{
		giveMoneySvc: giveMoneySvc,
	}
}

// GetAllRecords handles getting all give money records with pagination
// @Summary      获取所有打赏记录（分页）
// @Description  获取所有打赏记录，无需认证，支持分页
// @Tags         打赏管理
// @Produce      json
// @Param        page     query  int  false  "页码"  default(1)
// @Param        pageSize query  int  false  "每页数量"  default(10)
// @Success      200  {object}  response.Response{data=object{list=[]object{id=integer,nickname=string,figure=string,created_at=string,updated_at=string},total=integer,pageNum=integer,pageSize=integer}}  "获取成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /public/givemoney [get]
func (h *GiveMoneyHandler) GetAllRecords(c *gin.Context) {
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
	pageResult, err := h.giveMoneySvc.GetRecordsByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取打赏记录失败: "+err.Error())
		return
	}
	
	response.Success(c, gin.H{
		"list":     pageResult.Items,
		"total":    pageResult.Total,
		"pageNum":  page,
		"pageSize": pageSize,
	}, "获取打赏记录成功")
}

// CreateRecord handles creating a new give money record
// @Summary      创建打赏记录
// @Description  创建新的打赏记录，需要认证
// @Tags         打赏管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  object{nickname=string,figure=int}  true  "打赏记录信息"
// @Success      200  {object}  response.Response{data=object{id=integer,nickname=string,figure=int,created_at=string,updated_at=string}}  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "创建失败"
// @Router       /givemoney [post]
func (h *GiveMoneyHandler) CreateRecord(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname" binding:"required"`
		Figure   int    `json:"figure" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	record, err := h.giveMoneySvc.CreateRecord(c.Request.Context(), givemoney.CreateGiveMoneyParams{
		Nickname: req.Nickname,
		Figure:   req.Figure,
	})
	
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建打赏记录失败: "+err.Error())
		return
	}
	
	response.Success(c, record, "创建打赏记录成功")
}

// UpdateRecord handles updating an existing give money record
// @Summary      更新打赏记录
// @Description  更新打赏记录，需要认证
// @Tags         打赏管理
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int     true  "记录ID"
// @Param        body  body  object{nickname=string,figure=int}  true  "打赏记录信息"
// @Success      200  {object}  response.Response{data=object{id=integer,nickname=string,figure=int,created_at=string,updated_at=string}}  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误或ID非法"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "记录不存在"
// @Failure      500  {object}  response.Response  "更新失败"
// @Router       /givemoney/{id} [put]
func (h *GiveMoneyHandler) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
		return
	}
	
	var req struct {
		Nickname string `json:"nickname" binding:"required"`
		Figure   int    `json:"figure" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	
	record, err := h.giveMoneySvc.UpdateRecord(c.Request.Context(), uint(id), givemoney.UpdateGiveMoneyParams{
		Nickname: req.Nickname,
		Figure:   req.Figure,
	})
	
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新打赏记录失败: "+err.Error())
		return
	}
	
	response.Success(c, record, "更新打赏记录成功")
}

// DeleteRecord handles deleting a give money record
// @Summary      删除打赏记录
// @Description  删除打赏记录，需要认证
// @Tags         打赏管理
// @Security     BearerAuth
// @Param        id  path  int  true  "记录ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "ID非法"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "记录不存在"
// @Failure      500  {object}  response.Response  "删除失败"
// @Router       /givemoney/{id} [delete]
func (h *GiveMoneyHandler) DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
		return
	}
	
	if err := h.giveMoneySvc.DeleteRecord(c.Request.Context(), uint(id)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除打赏记录失败: "+err.Error())
		return
	}
	
	response.Success(c, nil, "删除打赏记录成功")
}