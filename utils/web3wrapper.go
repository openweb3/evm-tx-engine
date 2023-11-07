package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/web3go"
)

type TaggedBlockNumbers struct {
	LatestBlock    uint64
	SafeBlock      uint64
	FinalizedBlock uint64
}

type TxResultMeta struct {
	BlockNumber *hexutil.Big // which block transaction is in. nil if transaction is not contained in any block
	Status      uint8        // status field of transaction receipt
}

// NOTE: web3 wrapper should not rely any other package to prevent cycle refer

type Web3Wrapper struct {
}

var W3 Web3Wrapper

// TODO: return wrapper for different chain
func GetWeb3Client(chain string) (*web3go.Client, error) {
	return web3go.NewClient("http://localhost:8545")
}

// return latest, safe, finalized block number
func (w3 *Web3Wrapper) GetTaggedBlockNumbers(chain string) (TaggedBlockNumbers, error) {
	return TaggedBlockNumbers{
		30, 20, 10,
	}, nil
}

// nil if transaction is not on chain
// Returns transaction block number and status field
func (w3 *Web3Wrapper) GetTransactionResult(chain string, txHash string) (TxResultMeta, error) {
	blockNumber := big.NewInt(0)

	return TxResultMeta{
		(*hexutil.Big)(blockNumber), 0,
	}, nil
}
