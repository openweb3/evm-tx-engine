package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openweb3/evm-tx-engine/models"
	"gorm.io/gorm"
)

func initGin() *gin.Engine {
	engine := gin.New()
	return engine
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.Field{})
	db.AutoMigrate(&models.ChainTransaction{})
}

func main() {
	initGin()
	fmt.Println("Hello, World!")
}
