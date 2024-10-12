package controller

import (
	"github.com/gin-gonic/gin"
	"wallet/internal/apiResponse"
	"wallet/internal/message"
	"wallet/internal/requests"
	"wallet/internal/service"
	"wallet/pkg/hdwallet"
)

func NewMnemonic(c *gin.Context) {
	mnemonic, err := hdwallet.NewMnemonic()
	if err != nil {
		apiResponse.Fail(message.MnemonicCreateFail).Json(c)
		return
	}
	apiResponse.Success(mnemonic, message.Success).Json(c)
	return
}

func GetAddress(c *gin.Context) {
	var req requests.GetAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.Fail(message.ParamError).Json(c)
		return
	}
	service.GetWalletAddress(c, req).Json(c)
	return
}

func TransferUSDT(c *gin.Context) {
	var req requests.TransferUSDTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.Fail(message.ParamError).Json(c)
		return
	}
	service.TransferUSDT(c, req).Json(c)
	return
}

func GetTrxBalance(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		apiResponse.Fail(message.ParamError).Json(c)
		return
	}
	service.GetTrxBalance(c, address).Json(c)
	return
}

func GetUsdtBalance(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		apiResponse.Fail(message.ParamError).Json(c)
		return
	}
	service.GetUsdtBalance(c, address).Json(c)
	return
}

func GetEthUsdtBalance(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		apiResponse.Fail(message.ParamError).Json(c)
		return
	}
	service.GetEThUsdtBalance(c, address).Json(c)
	return
}
