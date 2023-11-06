package tests

import (
	"testing"
	"time"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/routers"
	"github.com/openweb3/evm-tx-engine/services"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initDbTestChain(db *gorm.DB) (*models.Chain, error) {
	chain := models.Chain{
		Name:           "ethereum",
		LatestBlock:    5,
		SafeBlock:      4,
		FinalizedBlock: 3,
	}
	err := db.First(&chain).Error
	if err == gorm.ErrRecordNotFound {
		err = db.Create(&chain).Error
	}
	return &chain, err
}

func initDbTestAccount(db *gorm.DB, chain *models.Chain) (*models.Account, error) {
	address := "0x000000000000000000000000"
	account := models.Account{
		Address: address,
	}
	err := db.Preload("ChainAccounts").First(&account).Error
	if err == gorm.ErrRecordNotFound || account.ChainAccounts == nil {
		account.ChainAccounts = []models.ChainAccount{
			{
				Chain:               *chain,
				Address:             address,
				LastSponsorInit:     time.Now(),
				LastSponsorReceived: time.Now(),
			},
		}
		err = db.Save(&account).Error
	}
	return &account, err
}

func insertTasks(db *gorm.DB) (*[]uint, error) {
	request := routers.TaskRequest{
		Chain: "ethereum",
		From:  "0x000000000000000000000000",
		Fields: routers.Fields{
			To: "0x000000000000000000000001",
		},
	}
	tasks := []uint{}
	for i := 0; i < 10; i++ {
		taskId, err := routers.CreateNewTask(db, request)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, taskId)
	}
	return &tasks, nil
}

func checkTasksStatus(db *gorm.DB, taskIds *[]uint) error {
	var tasks []models.Task
	err := db.Preload("History").Where("id IN ?", *taskIds).Find(&tasks).Error
	if err != nil {
		return err
	}
	for _, task := range tasks {
		logrus.WithField("taskId", task.ID).WithField("history", task.History)
	}
	return nil
}

func TestPipeline(t *testing.T) {
	db := models.ConnectDB("enginetest")
	chain, err := initDbTestChain(db)
	if err != nil {
		panic(err)
	}
	_, err = initDbTestAccount(db, chain)
	if err != nil {
		panic(err)
	}
	ctx, err := services.InitQueueContext(db)
	if err != nil {
		panic(err)
	}

	go services.StartWorkers(ctx)
	taskIds, err := insertTasks(db)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	err = checkTasksStatus(db, taskIds)
	if err != nil {
		panic(err)
	}
}
