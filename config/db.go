package config

import (
	"fmt"
	"go-chain-data/config/setting"
	"go-chain-data/global"
	models "go-chain-data/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBEngine(dbConfig *setting.DbConfig) (*gorm.DB, error) {
	conn := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(conn, dbConfig.Username, dbConfig.Pwd, dbConfig.Host, dbConfig.DbName, dbConfig.Charset, dbConfig.ParseTime)
	//fmt.Println("dsn", dsn)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MigrateDb 初始化数据库表
func MigrateDb() error {
	//AutoMigrate()，可以自动根据model创建/更新表结构
	if err := global.DBEngine.AutoMigrate(&models.Blocks{}, &models.Transaction{}, &models.Events{}, &models.Topic{}); err != nil {
		return err
	}
	return nil
}
