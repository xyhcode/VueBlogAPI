package router

import (
	"github.com/gin-gonic/gin"
)

// registerGiveMoneyRoutes 注册打赏记录相关的路由
func (r *Router) registerGiveMoneyRoutes(api *gin.RouterGroup) {
	// --- 前台公开接口 ---
	giveMoneyPublic := api.Group("/public/givemoney")
	{
		// 获取所有打赏记录: GET /api/public/givemoney
		giveMoneyPublic.GET("", r.giveMoneyHandler.GetAllRecords)
	}

	// --- 后台管理接口 ---
	giveMoneyAdmin := api.Group("/givemoney").Use(r.mw.JWTAuth(), r.mw.AdminAuth())
	{
		// 创建打赏记录: POST /api/givemoney
		giveMoneyAdmin.POST("", r.giveMoneyHandler.CreateRecord)
		
		// 更新打赏记录: PUT /api/givemoney/:id
		giveMoneyAdmin.PUT("/:id", r.giveMoneyHandler.UpdateRecord)
		
		// 删除打赏记录: DELETE /api/givemoney/:id
		giveMoneyAdmin.DELETE("/:id", r.giveMoneyHandler.DeleteRecord)
	}
}