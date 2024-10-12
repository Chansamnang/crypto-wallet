package grpcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
	"strings"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
)

type Client struct {
	node string
	GRPC *client.GrpcClient
}

func NewClient(node string) (*Client, error) {
	c := new(Client)
	c.node = node

	c.GRPC = client.NewGrpcClient(node)
	c.GRPC.SetTimeout(time.Second * 60)
	err := c.GRPC.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) SetTimeout(timeout time.Duration) error {
	if c == nil {
		return errors.New("client is nil ptr")
	}
	c.GRPC = client.NewGrpcClientWithTimeout(c.node, timeout)
	err := c.GRPC.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) keepConnect() error {
	_, err := c.GRPC.GetNodeInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return c.GRPC.Reconnect(c.node)
		}
		return err
	}
	return nil
}

func (c *Client) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	fmt.Printf("keepConnect %s", err)
	if err != nil {
		return nil, err
	}
	return c.GRPC.Transfer(from, to, amount)
}

func (c *Client) GetTrc10Balance(addr, assetId string) (int64, error) {
	err := c.keepConnect()
	if err != nil {
		return 0, err
	}
	acc, err := c.GRPC.GetAccount(addr)
	if err != nil || acc == nil {
		return 0, err
	}
	for key, value := range acc.AssetV2 {
		if key == assetId {
			return value, nil
		}
	}
	return 0, fmt.Errorf("%s do not find this assetID=%s amount", addr, assetId)
}

func (c *Client) GetTrxBalance(addr string) (*account.Account, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.GetAccountDetailed(addr)
}
func (c *Client) GetTrc20Balance(addr, contractAddress string) (*big.Int, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}

	return c.GRPC.TRC20ContractBalance(addr, contractAddress)
}

func (c *Client) TransferTrc10(from, to, assetId string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	fromAddr, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, errors.New("from address is not equal")
	}
	toAddr, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, errors.New("to address is not equal")
	}

	return c.GRPC.TransferAsset(fromAddr.String(), toAddr.String(), assetId, amount)
}

func (c *Client) TransferTrc20(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.TRC20Send(from, to, contract, amount, feeLimit)
}

func (c *Client) BroadcastTransaction(transaction *core.Transaction) error {
	err := c.keepConnect()
	if err != nil {
		return err
	}
	result, err := c.GRPC.Broadcast(transaction)
	if err != nil {
		return err
	}
	if result.Code != api.Return_SUCCESS {
		return errors.New("bad transaction: " + string(result.GetMessage()))
	}
	if result.Result == true {
		return nil
	}
	d, _ := json.Marshal(result)
	return errors.New("tx send fail: " + string(d))
}
