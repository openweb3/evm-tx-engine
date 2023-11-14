package models_test

import (
	"log"
	"testing"

	"github.com/holiman/uint256"
	"github.com/openweb3/evm-tx-engine/models"
	"gorm.io/gorm"
)

var db *gorm.DB

type ExampleModel struct {
	gorm.Model
	BigNumber *uint256.Int `gorm:"default:NULL"`
}

func init() {
	db = models.ConnectDB("enginetest")
	if err := db.AutoMigrate(&ExampleModel{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func TestUint256_SaveAndRetrieve(t *testing.T) {
	// 创建一个大整数值，用于测试
	bigIntValue := "1234567890123456789012345678901234567890123456789012345678901234567890"

	// 创建一个ExampleModel实例，并保存到数据库
	model := ExampleModel{BigNumber: uint256.MustFromDecimal(bigIntValue)}
	result := db.Create(&model)
	if result.Error != nil {
		t.Fatalf("Failed to save model: %v", result.Error)
	}

	// 从数据库检索模型
	var retrievedModel ExampleModel
	result = db.First(&retrievedModel, model.ID)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve model: %v", result.Error)
	}

	// 比较保存的值和检索的值
	if retrievedModel.BigNumber.Cmp(uint256.MustFromDecimal(bigIntValue)) != 0 {
		t.Errorf("Retrieved value %v does not match saved value %v", retrievedModel.BigNumber.String(), bigIntValue)
	}
}

func TestUint256_Default(t *testing.T) {
	// 创建一个大整数值，用于测试
	// 创建一个ExampleModel实例，并保存到数据库
	model := ExampleModel{}
	result := db.Create(&model)
	if result.Error != nil {
		t.Fatalf("Failed to save model: %v", result.Error)
	}

	// 从数据库检索模型
	var retrievedModel ExampleModel
	result = db.First(&retrievedModel, model.ID)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve model: %v", result.Error)
	}
}

// func TestUint256WithGorm(t *testing.T) {

// 	// 创建一个新记录...
// 	bigNum, _ := uint256.FromDecimal("12345678901234567890123456789012345678901234567890123456789012345678")
// 	myModel := ExampleModel{BigNumber: bigNum}
// 	result := db.Create(&myModel)
// 	if result.Error != nil {
// 		t.Fatalf("failed to create model: %v", result.Error)
// 	}

// 	// 从数据库检索记录...
// 	var retrievedModel ExampleModel
// 	db.First(&retrievedModel, "big_number = ?", bigNum.String())
// 	if retrievedModel.BigNumber.Cmp(bigNum) != 0 {
// 		t.Fatalf("retrieved number does not match: got %v want %v", retrievedModel.BigNumber, bigNum)
// 	}
// }
