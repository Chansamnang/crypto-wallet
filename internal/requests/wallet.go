package requests

import "github.com/shopspring/decimal"

type GetAddressRequest struct {
	Mnemonic string `json:"mnemonic" binding:"required"`
	Network  int    `json:"network" binding:"required"`
}

type TransferUSDTRequest struct {
	Mnemonic        string          `json:"mnemonic" binding:"required"`
	Network         int             `json:"network" binding:"required"`
	ReceiverAddress string          `json:"receiver_address" binding:"required"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
}
