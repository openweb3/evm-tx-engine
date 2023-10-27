package services

import (
	"fmt"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Scan GasEnoughQueue, and modify the gas field, then change the transaction to constructed (all fields are supposed to complete)

func StartPriceManageRound(db *gorm.DB) {
	var chainTransactions []models.ChainTransaction
	err := db.Joins("Field", "nonce IS not NULL AND gas_price = ?", 0).Where("tx_status = ?", utils.TxInternalGasEnoughQueue).Find(&chainTransactions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		} else {
			logrus.WithError(err).Error("failed to find tx to fill gas price")
		}
	}
	// ignore the case that gas price was already filled in the simple implememtation
	if len(chainTransactions) == 0 {
		return
	}
	for i := range chainTransactions {
		dbTransaction := db.Begin()
		err := func() error {
			if chainTransactions[i].Field.GasPrice != 0 || chainTransactions[i].Field.Nonce == nil {
				// TODO: check why this branch is reachable

				return fmt.Errorf("unexpected branch reached: gas price %d / nonce %d", chainTransactions[i].Field.GasPrice, chainTransactions[i].Field.Nonce)
			}
			chainTransactions[i].Field.GasPrice = 5
			chainTransactions[i].TxStatus = utils.TxInternalConstructed
			err := dbTransaction.Save(&(chainTransactions[i].Field)).Error
			if err != nil {
				return err
			}

			err = dbTransaction.Save(&(chainTransactions[i])).Error
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			logrus.WithError(err).WithField("txId", chainTransactions[i].ID).Error("Failed to update gas price")
			dbTransaction.Rollback()
			continue
		}
		err = dbTransaction.Commit().Error
		if err != nil {
			logrus.WithError(err).Error("Failed to commit gas price filling transaction")
			continue
		}

		logrus.WithField("txId", chainTransactions[i].ID).WithField("fieldId", chainTransactions[i].Field.ID).Info("Gas price filled")
	}

}
