package router

import (
	"github.com/gin-gonic/gin"
	"wallet/internal/controller"
)

func InitRouter(router gin.IRouter) {
	route := router.Group("/api/wallet")
	{
		route.GET("/newMnemonic", controller.NewMnemonic)
		route.POST("/address", controller.GetAddress)
		route.POST("/transferUsdt", controller.TransferUSDT)
		route.GET("/trxBalance", controller.GetTrxBalance)
		route.GET("/usdtBalance", controller.GetUsdtBalance)

		eth := route.Group("/eth")
		eth.GET("/usdtBalance", controller.GetEthUsdtBalance)
	}
}
