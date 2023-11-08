package accountadapter

import "github.com/openweb3/evm-tx-engine/models"

var Signer AccountAdapter

// AccountAdapter is the interface for the account manager adapter.
type AccountAdapter interface {
	// GetWrappingPublicKey returns the public key to wrap the imported secret key
	GetWrappingPublicKey() (map[string]string, error)
	// CreateAccount creates a new account with the given chain, wrapped private key, and wrapping public key
	CreateAccount(chain string, wrappedPrivateKey string, wrappingPublicKey string, wrappingAlg string) (string, error)
	// GetAccounts returns a list of accounts associated with the given chain
	GetAccounts(chain string) ([]models.Account, error)
	// SignRaw signs a raw byte array with the given address
	SignRaw(address string, data []byte) ([]byte, error)
	// SignTransaction signs a transaction with the given address
	SignTransaction(address string, transaction interface{}) ([]byte, error) // Transaction type will depend on specific blockchain
	// Encrypt encrypts the given data with the given address
	Encrypt(address string, data []byte) ([]byte, error)
	// Verify verifies the given signature and data with the given address
	Verify(address string, signature, data []byte) (bool, error)
}

func init() {
	Signer = NewSimpleAccountAdapter()
}
