package router

import (
	"github.com/gin-gonic/gin"
	"qtx/api/handler"
)

func Qtx(rg *gin.Engine) {
	// 测试功能
	group := rg.Group("qtx")
	{
		// 更新用户气泡状态
		group.POST("updateBubbleStatus", handler.UpdateBubbleStatus)
		// 更新用户兑换状态
		group.POST("updateExchangeStatus", handler.UpdateExchangeStatus)
	}
}
