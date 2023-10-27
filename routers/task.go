package routers

import (
	"fmt"

	"github.com/openweb3/evm-tx-engine/models"
	"gorm.io/gorm"
)

type TaskRequest struct {
	Chain     string `json:"chain"`
	From      string `json:"from"`
	Fields    Fields `json:"fields"`
	Fee       string `json:"fee"`
	RequestID string `json:"requestId"`
	Retry     Retry  `json:"retry"`
}

type Fields struct {
	To           string `json:"to"`
	Data         []byte `json:"data"`
	MaxFeePerGas uint   `json:"maxFeePerGas"`
	// Function  string        `json:"function"`
	// Params    []interface{} `json:"params"`
}

type Retry struct {
	MaxAttempts int64 `json:"maxAttempts"`
	Deadline    int64 `json:"deadline"`
}

// func SetupRoutes(router *gin.Engine) {
// 	apiV1 := router.Group("/v1")
// 	apiV1.POST("/task", CreateTask)
// }

// func CreateTask(c *gin.Context) {
// 	var taskRequest TaskRequest

// 	if err := c.ShouldBindJSON(&taskRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 这里假设我们有一个函数generateTaskID，它能生成一个新的任务ID
// 	// 以及一个函数createNewTask，它能根据给定的信息创建一个新任务
// 	taskID, err := createNewTask(taskRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, taskID)
// }

// 假设的函数来创建新任务，并返回任务ID
func CreateNewTask(db *gorm.DB, taskRequest TaskRequest) (uint, error) {
	// ... 在这里处理任务请求，创建新任务，并返回任务ID
	// var chain models.Chain
	// err := db.Where("name = ?", taskRequest.Chain).First(&chain).Error
	var chainAccount models.ChainAccount
	err := db.Preload("Chain", "name = ?", taskRequest.Chain).Where("address = ?", taskRequest.From).First(&chainAccount).Error
	if err != nil {
		return 0, fmt.Errorf("chain account %s %s not found", taskRequest.Chain, taskRequest.From)
	}
	task := models.Task{
		From: chainAccount,
		Field: models.Field{
			To:           taskRequest.Fields.To,
			Data:         taskRequest.Fields.Data,
			MaxFeePerGas: taskRequest.Fields.MaxFeePerGas,
		},
	}
	err = db.Create(&task).Error
	if err != nil {
		return 0, fmt.Errorf("create task error: %s", err.Error())
	}
	return task.ID, nil
}
