package accountadapter

import "github.com/openweb3/evm-tx-engine/models"

var Signer AccountAdapter

// AccountAdapter is the interface for the account manager adapter.
type AccountAdapter interface {
	GetWrappingPublicKey() (map[string]string, error)
	CreateAccount(chain string, wrappedPrivateKey string, wrappingPublicKey string, wrappingAlg string) (string, error)
	GetAccounts(chain string) ([]models.Account, error)
	SignRaw(address string, data []byte) ([]byte, error)
	SignTransaction(address string, transaction interface{}) ([]byte, error) // Transaction type will depend on specific blockchain
	Encrypt(address string, data []byte) ([]byte, error)
	Verify(address string, signature, data []byte) (bool, error)
}

func init() {
	Signer = NewSimpleAccountAdapter()
}
