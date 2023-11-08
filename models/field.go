package models

type Field struct {
	To string `gorm:"type:varchar(255)"`
	// Function      string             `gorm:"type:varchar(255)"`
	// Params        []string           `gorm:"type:json"`
	MaxFeePerGas      uint   `gorm:"type:int"`
	Data              []byte `gorm:"type:TINYBLOB"`
	Nonce             *uint  `gorm:"type:int"`
	GasLimit          uint   `gorm:"type:int"`
	GasPrice          uint   `gorm:"type:int"`
	Value             uint   `gorm:"type:int"`
	PriorityFeePerGas uint   `gorm:"type:int"`
}
