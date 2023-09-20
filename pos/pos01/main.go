package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//实现Pos
//定义区块

type Block struct {
	Index     int
	TimeStamp string
	BPM       int
	HashCode  string
	PreHash   string
	// 区块验证者
	Validator string
}

// 创建区块链，数组

var Blockchain []Block

// 生成新的区块
// address是矿工地址
// Block的BPM数据- BPM是每个验证者的人体脉搏值
func GenerateNextBlock(oldBlock Block, BPM int, address string) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = time.Now().String()
	newBlock.PreHash = oldBlock.HashCode
	newBlock.BPM = BPM
	newBlock.Validator = address
	newBlock.HashCode = GenerateHashValue(newBlock)
	return newBlock
}

// 哈希计算
func GenerateHashValue(block Block) string {
	var hashCode = block.PreHash + block.TimeStamp + block.Validator + strconv.Itoa(block.BPM) + strconv.Itoa(block.Index)
	var sha = sha256.New()
	sha.Write([]byte(hashCode))
	hashed := sha.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 网络上的全节点
type Node struct {
	// 记录有多少个币
	tokens int
	// 节点地址
	address string
}

// 存放几个节点，有几个用户在参与
var n [2]Node

// 用于记录挖矿地址
var addr [6000]string

func main() {
	// 测试
	//var firstBlock Block
	//myblock := GenerateNextBlock(firstBlock, 1, "adb")
	//fmt.Println(myblock)

	//创建两个参与者
	//传入持有的币和节点地址
	n[0] = Node{
		tokens:  1000,
		address: "abc123",
	}
	n[1] = Node{
		tokens:  5000,
		address: "bcd321",
	}

	// 以下是Pos共识算法
	var count = 0
	for i := 0; i < len(n); i++ {
		for j := 0; j < n[i].tokens; j++ {
			addr[count] = n[i].address
			count++
		}
	}

	// 设置随机种子
	//rand.Seed()
	rand.New(rand.NewSource(time.Now().Unix()))
	var rd = rand.Intn(6000)
	var adds = addr[rd]

	// 创建创世区块
	var firstBlock Block
	firstBlock.BPM = 100
	firstBlock.PreHash = "0"
	firstBlock.TimeStamp = time.Now().String()
	firstBlock.Validator = "abc123"
	firstBlock.Index = 1
	firstBlock.HashCode = GenerateHashValue(firstBlock)

	// 将区块加到区块链
	Blockchain = append(Blockchain, firstBlock)

	// 第二个区块
	// 让adds加入到
	var secondBlock = GenerateNextBlock(firstBlock, 200, adds)
	Blockchain = append(Blockchain, secondBlock)

	fmt.Println(Blockchain)

}
