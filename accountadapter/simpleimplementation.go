package accountadapter

// TODO: block same data signing (if not using same nonce, then secret key would leak)

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type SimpleAccountAdapter struct {
	keystore *keystore.KeyStore
}

func NewSimpleAccountAdapter() *SimpleAccountAdapter {
	// You can change the directory path and other parameters as needed.
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	return &SimpleAccountAdapter{
		keystore: ks,
	}
}

func (s *SimpleAccountAdapter) GetWrappingPublicKey() (map[string]string, error) {
	// Mock response, as this doesn't directly relate to Ethereum.
	return map[string]string{
		"wrappingPublicKey": "fnewoafnawoge",
		"wrappingAlg":       "fjweognawiowa",
	}, nil
}

func (s *SimpleAccountAdapter) CreateAccount(chain string, wrappedPrivateKey string, wrappingPublicKey string, wrappingAlg string) (string, error) {
	if chain != "ethereum" {
		return "", errors.New("unsupported chain")
	}

	// If a wrappedPrivateKey is provided, decrypt it using the wrappingPublicKey and wrappingAlg (this part is mocked for now)
	var privateKey *ecdsa.PrivateKey
	var err error
	if wrappedPrivateKey != "" {
		// Decrypt the privateKey (this is mocked for now)
		privateKey, err = crypto.HexToECDSA(wrappedPrivateKey)
		if err != nil {
			return "", err
		}
	} else {
		// If no privateKey is provided, create a new one.
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			return "", err
		}
	}

	// Import or create the account in the keystore.
	account, err := s.keystore.ImportECDSA(privateKey, "password") // You should have a secure way to handle passwords.
	if err != nil {
		return "", err
	}

	return account.Address.Hex(), nil
}

func (s *SimpleAccountAdapter) GetAccounts(chain string) ([]Account, error) {
	if chain != "ethereum" {
		return nil, errors.New("unsupported chain")
	}

	ethAccounts := s.keystore.Accounts()
	var accounts []Account
	for _, acc := range ethAccounts {
		accounts = append(accounts, Account{
			Address: acc.Address.Hex(),
			// ... other fields, you might want to store additional metadata separately.
		})
	}
	return accounts, nil
}

func (s *SimpleAccountAdapter) SignRaw(address string, data []byte) ([]byte, error) {
	account := accounts.Account{Address: common.HexToAddress(address)}
	signature, err := s.keystore.SignHash(account, crypto.Keccak256(data))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (s *SimpleAccountAdapter) SignTransaction(address string, transaction interface{}) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

func (s *SimpleAccountAdapter) Encrypt(address string, data []byte) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

func (s *SimpleAccountAdapter) Verify(address string, signature, data []byte) (bool, error) {
	return false, errors.New("Not implemented")
}
