package tron

import (
	"encoding/hex"
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"math/big"
	"wallet/pkg/common/config"
	grpcs "wallet/pkg/hdwallet/tron/grpc"
	"wallet/pkg/hdwallet/tron/sign"
	"wallet/pkg/zlogger"
)

var Client *grpcs.Client

func Init() {
	var err error
	Client, err = grpcs.NewClient(config.Config.Blockchain.TronGrpc)
	if err != nil {
		zlogger.Errorf("Failed to initialize TRON gRPC client: %v", err)
	}
	zlogger.Info("Initialized TRON gRPC client successfully")
}

func TransferTrc20(senderAddress string, receiverAddress string, fromPrivateKey []byte, amountTransfer *big.Int) (txid string, err error) {
	tx, err := Client.TransferTrc20(senderAddress, receiverAddress, config.Config.Blockchain.TronUSDTContract, amountTransfer, 50000000)
	if err != nil {
		return "", err
	}
	signTx, err := sign.TransactionSign(tx.Transaction, fromPrivateKey)
	if err != nil {
		return "", err
	}
	err = Client.BroadcastTransaction(signTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Txid), nil
}

func GetTrc20Balance(address string, contract string) (*big.Int, error) {
	return Client.GetTrc20Balance(address, contract)
}

func TransferTrc10(senderAddress string, receiverAddress string, fromPrivateKey []byte, amountTransfer int64) (txId string, err error) {
	tx, err := Client.TransferTrc10(senderAddress, receiverAddress, "trx", amountTransfer)
	if err != nil {
		return "", err
	}
	signTx, err := sign.TransactionSign(tx.Transaction, fromPrivateKey)
	if err != nil {
		return "", err
	}
	err = Client.BroadcastTransaction(signTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Txid), nil
}

func TransferTrx(senderAddress string, receiverAddress string, fromPrivateKey []byte, amountTransfer int64) (txid string, err error) {
	tx, err := Client.Transfer(senderAddress, receiverAddress, amountTransfer)
	if err != nil {
		return "", err
	}
	signTx, err := sign.TransactionSign(tx.Transaction, fromPrivateKey)
	if err != nil {
		return "", err
	}
	err = Client.BroadcastTransaction(signTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx.Txid), nil
}

func GetTrxBalance(address string) (*account.Account, error) {
	return Client.GetTrxBalance(address)
}
