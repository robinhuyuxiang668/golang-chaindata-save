package main

import (
	"context"
	"fmt"
	"go-chain-data/config"
	"go-chain-data/global"
	_ "go-chain-data/global"
	models "go-chain-data/internal/model"
	"go-chain-data/pkg/blockchain"
	"log"
	"math/big"
)

func init() {
	config.SetupConfig()
	config.SetupDBEngine()
	err := config.MigrateDb()
	if err != nil {
		log.Panic("config.MigrateDb error : ", err)
	}
	config.SetupEthClient()
}

func main() {
	//test()

	blockchain.InitBlock()
	//blockchain.SyncTask()
}

func test() {
	log.Println(global.BlockChainConfig.RpcUrl)
	block := models.Blocks{
		BlockHeight:       1,
		BlockHash:         "hash",
		ParentHash:        "parentHash",
		LatestBlockHeight: 2,
	}
	err := block.Insert()
	if err != nil {
		log.Panic("block.Insert error : ", err)
	}

	blockNumber, err := global.EthRpcClient.BlockNumber(context.Background())
	if err != nil {
		log.Panic("EthRpcClient.BlockNumber error : ", err)
	}
	log.Println("blockNumber is : ", blockNumber)

	lastBlock, err := global.EthRpcClient.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Panic("EthRpcClient.BlockByNumber error : ", err)
	}
	fmt.Println("lastBlock is : ", lastBlock)
}
