package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/services"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

func initGin() *gin.Engine {
	engine := gin.New()
	return engine
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Chain{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.ChainAccount{})
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.Field{})
	db.AutoMigrate(&models.ChainTransaction{})
}

func StartServices(db *gorm.DB) {
	// 500ms as default round interval
	defaultInterval := 500
	go utils.StartService(db, uint(defaultInterval), services.StartPickerRound)

	go utils.StartService(db, uint(defaultInterval), services.StartExecutionSimulationRound)
	go utils.StartService(db, uint(defaultInterval), services.StartBalanceCheckRound)

	go utils.StartService(db, uint(defaultInterval), services.StartNonceManageRound)
	go utils.StartService(db, uint(defaultInterval), services.StartPriceManageRound)

	go utils.StartService(db, uint(defaultInterval), services.StartSigningRound)
	go utils.StartService(db, uint(defaultInterval), services.StartSenderRound)

	go utils.StartService(db, uint(defaultInterval), services.StartTaggedBlockNumberUpdateRound)
	go utils.StartService(db, uint(defaultInterval), services.StartTransactionChainStatusUpdateRound)
	go utils.StartService(db, uint(defaultInterval), services.StartPendingStatusMoveRound)

}

func main() {
	initGin()
	fmt.Println("Hello, World!")
}
