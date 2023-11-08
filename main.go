package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func initGin() *gin.Engine {
	engine := gin.New()
	return engine
}

func main() {
	initGin()
	fmt.Println("evm tx engine, start!")
}
