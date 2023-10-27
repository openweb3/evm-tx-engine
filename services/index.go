package services

import (
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

func StartServices(db *gorm.DB) {
	// 500ms as default round interval
	defaultInterval := 100
	go utils.StartService(db, uint(defaultInterval), StartPickerRound)

	go utils.StartService(db, uint(defaultInterval), StartExecutionSimulationRound)
	go utils.StartService(db, uint(defaultInterval), StartBalanceCheckRound)

	go utils.StartService(db, uint(defaultInterval), StartNonceManageRound)
	go utils.StartService(db, uint(defaultInterval), StartPriceManageRound)

	go utils.StartService(db, uint(defaultInterval), StartSigningRound)
	go utils.StartService(db, uint(defaultInterval), StartSenderRound)

	go utils.StartService(db, uint(defaultInterval), StartTaggedBlockNumberUpdateRound)
	go utils.StartService(db, uint(defaultInterval), StartTransactionChainStatusUpdateRound)
	go utils.StartService(db, uint(defaultInterval), StartPendingStatusMoveRound)

}
