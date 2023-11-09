package models_test

import (
	"log"
	"math/big"
	"testing"

	"github.com/openweb3/evm-tx-engine/models"
	"gorm.io/gorm"
)

// func TestUint256_Scan(t *testing.T) {
// 	// 测试用例：从数据库读取的字节序列
// 	testCases := []struct {
// 		name     string
// 		input    interface{}
// 		expected string // 期望的大整数字符串表示形式
// 		wantErr  bool   // 是否期望有错误发生
// 	}{
// 		{
// 			name:     "Valid bytes",
// 			input:    []byte{0x01, 0x00, 0x00, 0x00},
// 			expected: "16777216", // 注意：假设数据库中的字节序列是大端序
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "Valid int64",
// 			input:    int64(42),
// 			expected: "42",
// 			wantErr:  false,
// 		},
// 		{
// 			name:    "Invalid type",
// 			input:   "invalid",
// 			wantErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var u models.Uint256
// 			err := u.Scan(tc.input)
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("Uint256.Scan() error = %v, wantErr %v", err, tc.wantErr)
// 				return
// 			}
// 			if err == nil && u.String() != tc.expected {
// 				t.Errorf("Uint256.Scan() got = %v, want %v", u.String(), tc.expected)
// 			}
// 		})
// 	}
// }

// func TestUint256_Value(t *testing.T) {
// 	// 测试用例：将大整数转换为字节序列
// 	testCases := []struct {
// 		name     string
// 		input    string // 输入的大整数字符串表示形式
// 		expected []byte // 期望的字节序列
// 		wantErr  bool   // 是否期望有错误发生
// 	}{
// 		{
// 			name:     "Valid number",
// 			input:    "16777216",
// 			expected: []byte{0x01, 0x00, 0x00, 0x00},
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "Big number",
// 			input:    "115792089237316195423570985008687907853269984665640564039457584007913129639935",
// 			expected: big.NewInt(0).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10).Bytes(),
// 			wantErr:  false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var u models.Uint256
// 			_, success := u.SetString(tc.input, 10)
// 			if !success {
// 				t.Fatalf("Failed to set Uint256 with string: %s", tc.input)
// 			}

// 			got, err := u.Value()
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("Uint256.Value() error = %v, wantErr %v", err, tc.wantErr)
// 				return
// 			}
// 			if !tc.wantErr && !compareBytes(got.([]byte), tc.expected) {
// 				t.Errorf("Uint256.Value() got = %v, want %v", got, tc.expected)
// 			}
// 		})
// 	}
// }

// // 辅助函数：比较两个字节序列是否相等
// func compareBytes(a, b []byte) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i, v := range a {
// 		if v != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// 假设Uint256类型和ExampleModel模型已经定义。
// ...

var db *gorm.DB

type ExampleModel struct {
	gorm.Model
	BigNumber models.Uint256
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
	bigInt, success := new(big.Int).SetString(bigIntValue, 10)
	if !success {
		t.Fatalf("Failed to set big.Int value")
	}

	// 创建一个ExampleModel实例，并保存到数据库
	model := ExampleModel{BigNumber: models.Uint256{*bigInt}}
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
	if retrievedModel.BigNumber.String() != bigIntValue {
		t.Errorf("Retrieved value %v does not match saved value %v", retrievedModel.BigNumber.String(), bigIntValue)
	}
}
