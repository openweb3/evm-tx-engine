package services

import (
	"github.com/openweb3/evm-tx-engine/models"
)

// Scan GasEnoughQueue, and modify the gas field, then change the transaction to constructed (all fields are supposed to complete)

// anyhow will allocate gas price
// don't operate db here because won't save transaction to db after gas price allocating
func SetGasPrice(tx *models.ChainTransaction) error {
	err := tx.Field.GasPrice.SetFromDecimal("1000000000")
	if err != nil {
		return err
	}
	return nil
}
