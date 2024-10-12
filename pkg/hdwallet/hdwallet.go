package hdwallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	geth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tronSdk "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"wallet/pkg/zlogger"
)

func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		zlogger.Errorf("Failed to generate entropy: %v", err)
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		zlogger.Errorf("Failed to generate mnemonic: %v", err)
		return "", err
	}
	return mnemonic, nil
}

func GetMasterKeyByMnemonic(mnemonic string) (*bip32.Key, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		zlogger.Errorf("Invalid mnemonic")
		return nil, errors.New("mnemonic is invalid")
	}
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		zlogger.Errorf("Failed to get master key: %v", err)
		return nil, err
	}

	return masterKey, nil
}

func DeriveEthAddress(privateKey *ecdsa.PrivateKey) geth.Address {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	ethAddress := crypto.PubkeyToAddress(*publicKey)
	return ethAddress
}

func DeriveTronAddress(privateKey *ecdsa.PrivateKey) tronSdk.Address {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	tronAddress := tronSdk.PubkeyToAddress(*publicKey)
	return tronAddress
}

func DerivePrivateKey(masterKey *bip32.Key, coinType uint32) (*ecdsa.PrivateKey, error) {
	// Derivation path: m/44'/coinType'/0'/0/0
	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose: %v", err)
	}

	coin, err := purpose.NewChildKey(bip32.FirstHardenedChild + coinType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type: %v", err)
	}

	account, err := coin.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %v", err)
	}

	change, err := account.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change: %v", err)
	}

	addressKey, err := change.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(addressKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDSA: %v", err)
	}

	return privateKey, nil

}
