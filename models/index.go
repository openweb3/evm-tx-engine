package models

import (
	"fmt"

	"github.com/openweb3/evm-tx-engine/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(dbName string) *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	var err error
	config.Init()
	dbConfig := config.GetConfig().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	err = Migrate(db)

	if err != nil {
		panic(err)
	}
	return db
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Chain{},
		&Account{},
		&ChainAccount{},
		&Task{},
		&Field{},
		&ChainTransaction{},
	)
	return err
}
