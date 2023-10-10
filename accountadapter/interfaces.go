package accountadapter

// Account represents a blockchain account with its metadata.
type Account struct {
	Address   string   `json:"address"`
	PublicKey string   `json:"publicKey"`
	Chains    []string `json:"chains"`
	Alias     string   `json:"alias"`
	Status    string   `json:"status"`
	Alg       string   `json:"alg"`
}

// AccountManagerAdapter is the interface for the account manager adapter.
type AccountManagerAdapter interface {
	GetWrappingPublicKey() (map[string]string, error)
	CreateAccount(chain string, wrappedPrivateKey string, wrappingPublicKey string, wrappingAlg string) (string, error)
	GetAccounts(chain string) ([]Account, error)
	SignRaw(address string, data []byte) ([]byte, error)
	SignTransaction(address string, transaction interface{}) ([]byte, error) // Transaction type will depend on specific blockchain
	Encrypt(address string, data []byte) ([]byte, error)
	Verify(address string, signature, data []byte) (bool, error)
}
