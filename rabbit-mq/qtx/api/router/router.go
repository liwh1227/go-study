package router

import (
	"qtx/api/handler/travel"
	"qtx/api/handler/wallet"

	"github.com/gin-gonic/gin"
)

// CarbonIntegral 碳积分相关api接口
func CarbonIntegral(rg *gin.Engine) {
	// 测试功能
	group := rg.Group("carbonIntegral")
	{
		group.POST("metro", travel.Metro)
		// 注册用户钱包
		group.POST("registerWallet", wallet.RegisterWallet)
	}
}
