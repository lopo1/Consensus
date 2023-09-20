package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	PreHash   string
	TimeStamp string
	Diff      int
	Data      string
	Index     int
	Nonce     int
	HashCode  string
}

/**
 * 创世区块
 */

func GenerateFirstBlock(data string) Block {
	var firstblock Block
	firstblock.PreHash = "0"
	firstblock.TimeStamp = time.Now().String()
	firstblock.Diff = 4
	firstblock.Data = data
	firstblock.Index = 1
	firstblock.Nonce = 0
	firstblock.HashCode = GenerationHashValue(firstblock)
	return firstblock

}

/**
 * GenerationHashValue
 * 生成区块hash
 */
func GenerationHashValue(block Block) string {
	var hashdata = strconv.Itoa(block.Index) + strconv.Itoa(block.Nonce) + strconv.Itoa(block.Diff) + block.TimeStamp
	var sha = sha256.New()
	sha.Write([]byte(hashdata))
	hashed := sha.Sum(nil)
	return hex.EncodeToString(hashed)
}

func main() {
	// 创建创世区块
	var firstBlock = GenerateFirstBlock("创世区块")
	fmt.Println("firstBlock= ", firstBlock)
	GenerateNextBlock("第二区块", firstBlock)

}

// 产生新的区块
func GenerateNextBlock(data string, oldBlock Block) Block {
	var newBlock Block
	newBlock.TimeStamp = time.Now().String()
	newBlock.Diff = 4
	newBlock.Index = 2
	newBlock.Nonce = 0
	newBlock.PreHash = oldBlock.PreHash

	// 创建pow算法
	// 计算前导0为4个的hash值
	newBlock.HashCode = Pow(newBlock.Diff, &newBlock)
	return newBlock
}

func Pow(diff int, block *Block) string {
	for {
		hash := GenerationHashValue(*block)
		fmt.Println("powtest01 hash", hash)
		//strings.Repeat("0", diff) 判断hash是否有diff 个0
		if strings.HasPrefix(hash, strings.Repeat("0", diff)) {
			fmt.Println("挖矿成功")
			return hash
		} else {
			block.Nonce++
		}
	}
}
