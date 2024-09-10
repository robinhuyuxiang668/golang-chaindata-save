package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go-chain-data/global"
	models "go-chain-data/internal/model"
	"log"
	"math/big"
	"time"
)

// InitBlock 初始化第一个区块数据
/**
先查询数据库中是否已经存在block记录
若不存在，查询最新区块高度
通过区块高度查询最新区块信息
组装数据，存储到数据库
*/
func InitBlock() {
	block := &models.Blocks{}
	count := block.Counts()
	if count == 0 {
		lastBlockNumber, err := global.EthRpcClient.BlockNumber(context.Background())
		if err != nil {
			log.Panic("InitBlock - BlockNumber err : ", err)
		}
		lastBlock, err := global.EthRpcClient.BlockByNumber(context.Background(), big.NewInt(int64(lastBlockNumber)))

		if err != nil {
			log.Panic("InitBlock - BlockByNumber err : ", err)
		}
		block.BlockHash = lastBlock.Hash().Hex()
		block.BlockHeight = lastBlock.NumberU64()
		block.LatestBlockHeight = lastBlock.NumberU64()
		block.ParentHash = lastBlock.ParentHash().Hex()
		err = block.Insert()
		if err != nil {
			log.Panic("InitBlock - Insert block err : ", err)
		}
	}
}

// SyncTask
// 间隔1S从链上拉取数据
// 获取最新区块高度
// 从数据库查询最新存储的区块数据
// 判断数据库存储的最新区块链高度是否大于查询的最新区块高度
// 如果大于则跳出循环不执行后面操作，反之通过数据库存储的最新区块链高度查询区块信息
// 通过HandleBlock()方法处理最新区块信息（存储到数据库）
func SyncTask() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("start get latestBlockNumber...")
			latestBlockNumber, err := global.EthRpcClient.BlockNumber(context.Background())
			if err != nil {
				log.Panic("EthRpcClient.BlockNumber error : ", err)
			}
			var blocks models.Blocks
			latestBlock, err := blocks.GetLatest()
			if err != nil {
				log.Panic("blocks.GetLatest error : ", err)
			}
			if latestBlock.LatestBlockHeight > latestBlockNumber {
				log.Printf("databse's LatestBlockHeight : %v greater than current latestBlockNumber : %v \n", latestBlock.LatestBlockHeight, latestBlockNumber)
				continue
			}
			currentBlock, err := global.EthRpcClient.BlockByNumber(context.Background(), big.NewInt(int64(latestBlock.LatestBlockHeight)))
			if err != nil {
				log.Panic("EthRpcClient.BlockByNumber error : ", err)
			}
			log.Printf("get currentBlock blockNumber : %v , blockHash : %v \n", currentBlock.Number(), currentBlock.Hash().Hex())
			err = HandleBlock(currentBlock)
			if err != nil {
				log.Panic("HandleBlock error : ", err)
			}
		}
	}
}

// HandleBlock
// 处理区块数据，存储到数据库
// 调用HandleTransaction()方法处理区块里包含的交易数据
func HandleBlock(currentBlock *types.Block) error {
	block := &models.Blocks{
		BlockHeight:       currentBlock.NumberU64(),
		BlockHash:         currentBlock.Hash().Hex(),
		ParentHash:        currentBlock.ParentHash().Hex(),
		LatestBlockHeight: currentBlock.NumberU64() + 1,
	}
	err := block.Insert()
	if err != nil {
		return err
	}
	err = HandleTransaction(currentBlock)
	if err != nil {
		return err
	}
	return nil
}

// 判断一个地址是否是合约地址
func isContractAddress(address string) (bool, error) {

	addr := common.HexToAddress(address)
	code, err := global.EthRpcClient.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return false, err
	}

	return len(code) > 0, nil
}
