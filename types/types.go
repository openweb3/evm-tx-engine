package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/evm-tx-engine/types/code"
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

// why legacy tx status is needed
// certain transaction will be moved into error state because of too long pending()
func InferSentTransactionStatus(meta TxResultMeta, legacyTxStatus code.TxStatus, cachedTaggedBlockNumbers TaggedBlockNumbers) (code.TxStatus, error) {
	// TODO: consider errors
	if !legacyTxStatus.IsSent() {
		return legacyTxStatus,
			fmt.Errorf("transaction status is not sent, but it is sent")
	}
	if meta.BlockNumber == nil {
		return code.TxPoolPending, nil
	}
	if meta.BlockNumber.ToInt().Uint64() <= cachedTaggedBlockNumbers.FinalizedBlock {
		return code.TxChainFinalized, nil
	}
	if meta.BlockNumber.ToInt().Uint64() <= cachedTaggedBlockNumbers.SafeBlock {
		return code.TxChainSafe, nil
	}
	if meta.BlockNumber.ToInt().Uint64() <= cachedTaggedBlockNumbers.LatestBlock {
		return code.TxChainLatest, nil
	}
	return legacyTxStatus,
		fmt.Errorf("transaction not in expected status")
}
