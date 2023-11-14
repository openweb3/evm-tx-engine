package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/evm-tx-engine/types"
	"github.com/openweb3/web3go"
)

// NOTE: web3 wrapper should not rely any other package to prevent cycle refer

type Web3Wrapper struct {
}

var W3 Web3Wrapper

// TODO: return wrapper for different chain
func GetWeb3Client(chain string) (*web3go.Client, error) {
	return web3go.NewClient("http://localhost:8545")
}

// return latest, safe, finalized block number
func (w3 *Web3Wrapper) GetTaggedBlockNumbers(chain string) (types.TaggedBlockNumbers, error) {
	return types.TaggedBlockNumbers{
		LatestBlock: 30, SafeBlock: 20, FinalizedBlock: 10,
	}, nil
}

// nil if transaction is not on chain
// Returns transaction block number and status field
func (w3 *Web3Wrapper) GetTransactionResult(chain string, txHash string) (types.TxResultMeta, error) {
	blockNumber := big.NewInt(0)

	return types.TxResultMeta{
		BlockNumber: (*hexutil.Big)(blockNumber), Status: 0,
	}, nil
}
