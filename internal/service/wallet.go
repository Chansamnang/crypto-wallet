package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"regexp"
	"wallet/internal/apiResponse"
	"wallet/internal/message"
	"wallet/internal/requests"
	"wallet/pkg/common/config"
	"wallet/pkg/constant"
	"wallet/pkg/hdwallet"
	"wallet/pkg/hdwallet/eth"
	"wallet/pkg/hdwallet/tron"
	"wallet/pkg/tools/blockchain"
	"wallet/pkg/zlogger"
)

func GetWalletAddress(c context.Context, req requests.GetAddressRequest) *apiResponse.Response {
	var err error
	masterKey, err := hdwallet.GetMasterKeyByMnemonic(req.Mnemonic)
	if err != nil {
		return apiResponse.Fail(message.InvalidMnemonic)
	}

	switch req.Network {
	case constant.NetworkTron:
		privateKey, err := hdwallet.DerivePrivateKey(masterKey, constant.CoinTron)
		if err != nil {
			return apiResponse.Fail(message.Fail)
		}
		tronAddr := hdwallet.DeriveTronAddress(privateKey)
		return apiResponse.Success(tronAddr.String(), message.Success)
	case constant.NetworkEth:
		privateKey, err := hdwallet.DerivePrivateKey(masterKey, constant.CoinEth)
		if err != nil {
			return apiResponse.Fail(message.Fail)
		}
		ethAddr := hdwallet.DeriveEthAddress(privateKey)
		return apiResponse.Success(ethAddr.Hex(), message.Success)
	default:
		return apiResponse.Fail(message.ParamError)
	}
}

func TransferUSDT(ctx context.Context, req requests.TransferUSDTRequest) *apiResponse.Response {
	var err error
	masterKey, err := hdwallet.GetMasterKeyByMnemonic(req.Mnemonic)
	if err != nil {
		return apiResponse.Fail(message.Fail)
	}

	switch req.Network {
	case constant.NetworkTron:
		privateKey, err := hdwallet.DerivePrivateKey(masterKey, constant.CoinTron)
		if err != nil {
			return apiResponse.Fail(message.InvalidMnemonic)
		}
		senderAddr := hdwallet.DeriveTronAddress(privateKey)
		if senderAddr.String() == req.ReceiverAddress {
			return apiResponse.Fail(message.SelfTransferNotAllow)
		}

		availableUSDT, err := tron.GetTrc20Balance(senderAddr.String(), config.Config.Blockchain.TronUSDTContract)
		availableUsdt := blockchain.ToDecimal(availableUSDT, constant.UsdtDecimals)
		if err != nil {
			zlogger.Errorf("[TransferUSDT] GetTrc20Balance error %v", err)
			return apiResponse.Fail(message.Fail)
		}

		if availableUsdt.Cmp(req.Amount) < 0 {
			zlogger.Warnf("Available Balance sender %s, balance %v, transfer amount %v", senderAddr.String(), availableUsdt, req.Amount)
			return apiResponse.Fail(message.LowBalance)
		}

		txId, err := tron.TransferTrc20(senderAddr.String(), req.ReceiverAddress, crypto.FromECDSA(privateKey), blockchain.ToWei(req.Amount, constant.UsdtDecimals))
		if err != nil {
			zlogger.Errorf("[TransferUSDT] transfer fail %v", err)
			return apiResponse.Fail(message.Fail)
		}
		return apiResponse.Success(txId, message.Success)
	case constant.NetworkEth:
		privateKey, err := hdwallet.DerivePrivateKey(masterKey, constant.CoinEth)
		if err != nil {
			return apiResponse.Fail(message.InvalidMnemonic)
		}
		senderAddr := hdwallet.DeriveEthAddress(privateKey)
		if senderAddr.String() == req.ReceiverAddress {
			return apiResponse.Fail(message.SelfTransferNotAllow)
		}
		gasPrice, err := eth.Client.SuggestGasPrice(ctx)
		if err != nil {
			return apiResponse.Fail(message.Fail)
		}
		txId, err := eth.Client.TransferUSDT(ctx, crypto.FromECDSA(privateKey), senderAddr, req.ReceiverAddress, blockchain.ToWei(req.Amount, constant.UsdtDecimals), gasPrice)
		if err != nil {
			zlogger.Errorf("[TransferUSDT] transfer fail %v", err)
			return apiResponse.Fail(message.Fail)
		}
		return apiResponse.Success(txId, message.Success)
	default:
		return apiResponse.Fail(message.ParamError)
	}
}

func GetTrxBalance(ctx context.Context, address string) *apiResponse.Response {
	result, _ := validateAddressFmt(address, constant.NetworkTron)
	if !result {
		return apiResponse.Fail(message.InvalidAddressFormat)
	}
	account, err := tron.GetTrxBalance(address)
	if err != nil {
		zlogger.Errorf("[GetTrxBalance] get tron balance error %v", err)
		return apiResponse.Fail(message.Fail)
	}
	return apiResponse.Success(blockchain.ToDecimal(account.Balance, constant.UsdtDecimals), message.Success)
}

func GetUsdtBalance(ctx context.Context, address string) *apiResponse.Response {
	account, err := tron.GetTrc20Balance(address, config.Config.Blockchain.TronUSDTContract)
	if err != nil {
		zlogger.Errorf("[GetTrxBalance] get usdt balance error %v", err)
		return apiResponse.Fail(message.Fail)
	}
	return apiResponse.Success(blockchain.ToDecimal(account, constant.UsdtDecimals), message.Success)
}

func GetEThUsdtBalance(ctx context.Context, address string) *apiResponse.Response {
	account, err := eth.Client.GetUSDT(address)
	if err != nil {
		zlogger.Errorf("[GetEThUsdtBalance] get usdt balance error %v", err)
		return apiResponse.Fail(message.Fail)
	}
	return apiResponse.Success(blockchain.ToDecimal(account, constant.UsdtDecimals), message.Success)
}

func validateAddressFmt(address string, network int) (bool, error) {
	switch network {
	case constant.NetworkTron:
		result, err := regexp.MatchString(constant.TronAddressFmt, address)
		return result, err
	case constant.NetworkEth:
		result, err := regexp.MatchString(constant.EthAddressFmt, address)
		return result, err
	default:
		return false, fmt.Errorf("invalid network")
	}
}
