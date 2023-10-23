package utils

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/web3go"
)

type TaggedBlockNumbers struct {
	LatestBlock    uint64
	SafeBlock      uint64
	FinalizedBlock uint64
}

type TxStatusMeta struct {
	BlockNumber *hexutil.Big
	Status      uint64 // Status field of transaction receipt
}

// NOTE: web3 wrapper should not rely any other package to prevent cycle refer

type Web3Wrapper struct {
}

var W3 Web3Wrapper

// TODO: return wrapper for different chain
func GetWeb3Client(chain string) (*web3go.Client, error) {
	return web3go.NewClient("http://localhost:8545")
}

// func (w3 *Web3Wrapper) GetSingleChainTransactionStatus(txHash string) (TxStatus, error) {
// 	return TxPoolPending, nil
// }

// return latest, safe, finalized block number
func (w3 *Web3Wrapper) GetTaggedBlockNumbers(chain string) (TaggedBlockNumbers, error) {
	return TaggedBlockNumbers{
		10, 20, 30,
	}, nil
}

// nil if transaction is not on chain
// Returns transaction block number and status field
func (w3 *Web3Wrapper) GetTransactionDetail(chain string, txHash string) (TxStatusMeta, error) {
	return TxStatusMeta{
		nil, 0,
	}, nil
}
