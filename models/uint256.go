package models

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

// Uint256 is a type that holds a big integer, can be used to store uint256 types.
type Uint256 struct {
	big.Int
}

// Scan 实现了 sql.Scanner 接口，用于从数据库读取值。
func (u *Uint256) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		// 数据库中是以字节序列的形式存储大整数。
		u.SetInt64(0).SetBytes(v)
		return nil
	case int64:
		// 如果你知道数据库中的值是小于uint256的整数范围，可以直接使用。
		u.SetInt64(v)
		return nil
	default:
		return errors.New("unsupported type")
	}
}

// Value 实现了 driver.Valuer 接口，用于将值写入数据库。
func (u Uint256) Value() (driver.Value, error) {
	// 将大整数转换为字节序列。
	return u.Bytes(), nil
}
