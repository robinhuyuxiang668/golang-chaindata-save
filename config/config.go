package config

import (
	"github.com/spf13/viper"
	"go-chain-data/global"
	"log"
)

type Config struct {
	vp *viper.Viper
}

// NewConfig
// 读取yml配置信息并 创建viper的实例对象
func NewConfig() (*Config, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("config")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Config{vp}, nil
}

// ReadSection
// 通过给定的 k值 读取配置文件对应的配置信息并存到v变量
func (config *Config) ReadSection(k string, v interface{}) error {
	err := config.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

func SetupConfig() {
	conf, err := NewConfig()
	if err != nil {
		log.Panic("NewConfig error : ", err)
	}

	//global.DbConfig是指针类型,初始是nil,但是nil对应的地址却存在,
	//所以这里传&;如果直接传global.DbConfig即传nil报错因为根本找不到value
	err = conf.ReadSection("Database", &global.DbConfig)
	if err != nil {
		log.Panic("ReadSection - Database error : ", err)
	}
	err = conf.ReadSection("BlockChain", &global.BlockChainConfig)
	if err != nil {
		log.Panic("ReadSection - BlockChain error : ", err)
	}
}
func SetupDBEngine() {
	var err error
	global.DBEngine, err = NewDBEngine(global.DbConfig)
	if err != nil {
		log.Panic("NewDBEngine error : ", err)
	}
}

func SetupEthClient() {
	var err error
	global.EthRpcClient, err = NewEthRpcClient()
	if err != nil {
		log.Panic("NewEthRpcClient error : ", err)
	}
}
