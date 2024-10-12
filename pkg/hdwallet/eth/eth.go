package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"wallet/pkg/common/config"
	"wallet/pkg/zlogger"
)

type Geth interface {
	GetETH(ctx context.Context, address string) (*big.Float, error)
	GetUSDT(address string) (*big.Int, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGasLimit(ctx context.Context, address string, data []byte) (uint64, error)
	GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	GetTransactionInfo(ctx context.Context, hash common.Hash) (*types.Receipt, *types.Transaction, bool, error)
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
	GetLatestBlock(ctx context.Context) (*big.Int, error)
	GetTransactionByBlockNum(ctx context.Context, blockNum *big.Int) (types.Transactions, error)
	GetChainID(ctx context.Context) (*big.Int, error)
	TransferETH(
		ctx context.Context,
		senderPrivateKey []byte,
		senderAddress common.Address,
		receiverPublicKey string,
		amount *big.Int,
		gasPrice *big.Int,
	) (string, error)

	TransferUSDT(ctx context.Context,
		senderPrivateKey []byte,
		senderAddress common.Address,
		receiverPublicKey string,
		amount *big.Int,
		gasPrice *big.Int,
	) (string, error)
	ETHDecimals() int
	USDTDecimals() int
	LeftPadBytesLength() int
	Erc20GasLimit() uint64
	ETHGasLimit() uint64
	TransactionNotFoundMsg() string
}

type geth struct {
	client *ethclient.Client
}

const (
	ethDecimals         = 18
	usdtDecimals        = 6
	leftPadBytesLength  = 32
	erc20GasLimit       = 100000
	ethGasLimit         = 21000
	transactionNotFound = "not found"
)

var (
	ErrFailToParse = errors.New("fail to parse big float from string")
	ErrCastToECDSA = errors.New("fail to cast ECDSA")
)

var Client Geth

func Init() Geth {
	client, err := ethclient.Dial(config.Config.Blockchain.EthAlchemy)
	if err != nil {
		zlogger.Errorf("connect eth blockchain error %v", err)
		panic(err)
	}
	Client = &geth{client: client}
	return Client
}

func (g geth) GetETH(ctx context.Context, address string) (*big.Float, error) {
	account := common.HexToAddress(address)

	weiBalance, err := g.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return &big.Float{}, err
	}

	weiBalanceBigFloat, success := new(big.Float).SetString(weiBalance.String())
	if !success {
		return &big.Float{}, ErrFailToParse
	}

	balance := new(big.Float).Quo(weiBalanceBigFloat, big.NewFloat(math.Pow10(ethDecimals)))

	return balance, nil
}

func (g geth) GetUSDT(address string) (*big.Int, error) {
	instance, err := NewToken(common.HexToAddress(config.Config.Blockchain.EthUSDTContract), g.client)
	if err != nil {
		return &big.Int{}, err
	}

	weiBalance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return &big.Int{}, err
	}

	return weiBalance, nil
}

func (g geth) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := g.client.SuggestGasPrice(ctx)
	if err != nil {
		return &big.Int{}, err
	}

	return gasPrice, nil
}

func (g geth) EstimateGasLimit(ctx context.Context, address string, data []byte) (uint64, error) {
	a := common.HexToAddress(address)

	gasLimit, err := g.client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &a,
		Data: data,
	})
	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}

func (g geth) TransferETH(
	ctx context.Context,
	senderPrivateKey []byte,
	senderAddress common.Address,
	receiverPublicKey string,
	amount *big.Int,
	gasPrice *big.Int,
) (string, error) {
	nonce, err := g.client.PendingNonceAt(ctx, senderAddress)
	if err != nil {
		return "", err
	}

	chainID, err := g.GetChainID(ctx)
	if err != nil {
		return "", err
	}

	userAddressPk := common.HexToAddress(receiverPublicKey)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &userAddressPk,
		Value:    amount,
		Gas:      erc20GasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), crypto.ToECDSAUnsafe(senderPrivateKey))

	if err != nil {
		return "", err
	}

	err = g.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func (g geth) TransferUSDT(
	ctx context.Context,
	senderPrivateKey []byte,
	senderAddress common.Address,
	receiverPublicKey string,
	amount *big.Int,
	gasPrice *big.Int,
) (string, error) {
	if amount.Cmp(big.NewInt(0)) > 0 {
		nonce, err := g.client.PendingNonceAt(ctx, senderAddress)
		if err != nil {
			return "", err
		}

		receiverAddress := common.HexToAddress(receiverPublicKey)

		hash := sha3.NewLegacyKeccak256()
		hash.Write([]byte("transfer(address,uint256)"))
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(receiverAddress.Bytes(), leftPadBytesLength)
		paddedAmount := common.LeftPadBytes(amount.Bytes(), leftPadBytesLength)

		var data []byte
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)

		chainID, err := g.GetChainID(ctx)
		if err != nil {
			return "", err
		}

		contract := common.HexToAddress(config.Config.Blockchain.EthUSDTContract)

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &contract,
			Value:    big.NewInt(0),
			Gas:      erc20GasLimit,
			GasPrice: gasPrice,
			Data:     data,
		})

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), crypto.ToECDSAUnsafe(senderPrivateKey))
		if err != nil {
			return "", err
		}

		err = g.client.SendTransaction(ctx, signedTx)
		if err != nil {
			return "", err
		}

		return signedTx.Hash().Hex(), nil
	}

	return "", nil

}

func (g geth) GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := g.client.TransactionByHash(ctx, hash)
	if err != nil {
		return &types.Transaction{}, isPending, err
	}

	return tx, isPending, nil
}

func (g geth) GetTransactionInfo(ctx context.Context, hash common.Hash) (
	*types.Receipt,
	*types.Transaction,
	bool,
	error,
) {
	tx, isPending, err := g.GetTransaction(ctx, hash)
	if err != nil {
		return &types.Receipt{}, &types.Transaction{}, isPending, err
	}

	receipt, err := g.GetTransactionReceipt(ctx, hash)
	if err != nil {
		return &types.Receipt{}, &types.Transaction{}, isPending, err
	}

	return receipt, tx, isPending, nil
}

func (g geth) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	receipt, err := g.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return &types.Receipt{}, err
	}

	return receipt, nil
}

func (g geth) GetLatestBlock(ctx context.Context) (*big.Int, error) {
	header, err := g.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return &big.Int{}, err
	}

	return header.Number, nil
}

func (g geth) GetTransactionByBlockNum(ctx context.Context, blockNum *big.Int) (types.Transactions, error) {
	block, err := g.client.BlockByNumber(ctx, blockNum)
	if err != nil {
		return nil, err
	}

	return block.Transactions(), nil
}

func (g geth) GetChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := g.client.NetworkID(ctx)
	if err != nil {
		return &big.Int{}, err
	}

	return chainID, nil
}

// zeroKey zeroes a private key in memory.
func (g geth) zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func (g geth) ETHDecimals() int {
	return ethDecimals
}

func (g geth) USDTDecimals() int {
	return usdtDecimals
}

func (g geth) LeftPadBytesLength() int {
	return leftPadBytesLength
}

func (g geth) Erc20GasLimit() uint64 {
	return erc20GasLimit
}

func (g geth) ETHGasLimit() uint64 {
	return ethGasLimit
}

func (g geth) TransactionNotFoundMsg() string {
	return transactionNotFound
}
