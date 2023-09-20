package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Pos思想，按照持有币的比重，分配记账权的比重

// 定义全节点
type PNode struct {
	// 持有币的数量
	Tokens int
	//持币时间
	Days int
	// 地址
	Address string
}

// PBlock
type PBlock struct {
	Index     int
	Data      string
	PreHash   string
	Hash      string
	TimeStamp string
	// 区块验证者
	Validator *PNode
}

// FirstBlock 生成创世区块
func FirstBlock() PBlock {
	var firstBlock = PBlock{
		0,
		"创世区块",
		"",
		"",
		time.Now().String(),
		&PNode{
			Tokens:  0,
			Days:    0,
			Address: "",
		},
	}
	firstBlock.Hash = hex.EncodeToString(BlockHash(&firstBlock))
	return firstBlock
}

// 计算hash
func BlockHash(block *PBlock) []byte {
	hashed := strconv.Itoa(block.Index) + block.Data + block.PreHash +
		block.TimeStamp + block.Validator.Address
	h := sha256.New()
	h.Write([]byte(hashed))
	hash := h.Sum(nil)
	return hash
}

// 创建5个全节点
var nodes = make([]PNode, 5)

// 存放节点的地址
var addr = make([]*PNode, 15)

func InitNodes() {
	nodes[0] = PNode{
		Tokens:  1,
		Days:    1,
		Address: "0x123456",
	}
	nodes[1] = PNode{
		Tokens:  2,
		Days:    1,
		Address: "0x233456",
	}
	nodes[2] = PNode{
		Tokens:  3,
		Days:    1,
		Address: "0x343456",
	}
	nodes[3] = PNode{
		Tokens:  4,
		Days:    1,
		Address: "0x453456",
	}
	nodes[4] = PNode{
		Tokens:  5,
		Days:    1,
		Address: "0x563456",
	}

	count := 0
	for i := 0; i < len(nodes); i++ {
		// 持币数量* 币龄
		for j := 0; j < nodes[i].Tokens*nodes[i].Days; j++ {
			addr[count] = &nodes[i]
			count++
		}
	}
}

func CreateNewBlock(lastBlock *PBlock, data string) PBlock {
	// 生成区块
	var newBlock PBlock
	newBlock.Index = lastBlock.Index + 1
	newBlock.TimeStamp = time.Now().String()
	newBlock.PreHash = lastBlock.Hash
	newBlock.Data = data

	// 需要休眠一下
	time.Sleep(1 * time.Second)
	// 产生[0,15)的随机数
	rand.NewSource(time.Now().Unix())
	var rd = rand.Intn(15)
	// 选出矿工
	node := addr[rd]
	fmt.Printf("由%s 根据Pos算法生成新的区块\n", node.Address)
	//验证者，实际的挖矿人
	newBlock.Validator = node
	// 模拟挖矿所得奖励
	node.Tokens += 1
	newBlock.Hash = hex.EncodeToString(BlockHash(&newBlock))
	return newBlock

}

func main() {
	InitNodes()
	var firstBlock = FirstBlock()
	// 创建新的区块
	for i := 0; i < 30; i++ {
		var newBlock = CreateNewBlock(&firstBlock, "新的区块")
		fmt.Println("新的区块")
		fmt.Println(newBlock)
	}
}
