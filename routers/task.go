package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	Data         string `json:"data"`
	MaxFeePerGas string `json:"maxFeePerGas"`
	// Function  string        `json:"function"`
	// Params    []interface{} `json:"params"`
}

type Retry struct {
	MaxAttempts int64 `json:"maxAttempts"`
	Deadline    int64 `json:"deadline"`
}

func SetupRoutes(router *gin.Engine) {
	apiV1 := router.Group("/v1")
	apiV1.POST("/task", CreateTask)
}

func CreateTask(c *gin.Context) {
	var taskRequest TaskRequest

	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里假设我们有一个函数generateTaskID，它能生成一个新的任务ID
	// 以及一个函数createNewTask，它能根据给定的信息创建一个新任务
	taskID, err := createNewTask(taskRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskID)
}

// 假设的函数来创建新任务，并返回任务ID
func createNewTask(taskRequest TaskRequest) (string, error) {
	// ... 在这里处理任务请求，创建新任务，并返回任务ID
	return "1234556", nil
}
