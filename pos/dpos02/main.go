package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 定义全节点
type Node struct {
	// 节点名称
	Name string
	//被选举的票数
	Votes int
}

// Block 定义区块
type Block struct {
	Index     int
	PreHash   string
	Hash      string
	TimeStamp string
	Data      []byte
	// 代理人
	delegate *Node
}

func firstBlock() Block {
	generate := Block{
		Index:     0,
		PreHash:   "",
		Hash:      "",
		TimeStamp: time.Now().String(),
		Data:      []byte("first block"),
		delegate:  nil,
	}
	generate.Hash = hex.EncodeToString(blockHash(generate))
	return generate
}

// 计算hash
func blockHash(block Block) []byte {
	hash := block.PreHash + block.TimeStamp + block.delegate.Name + strconv.Itoa(block.delegate.Votes) + hex.EncodeToString(block.Data) + strconv.Itoa(block.Index)
	var sha = sha256.New()
	sha.Write([]byte(hash))
	hashed := sha.Sum(nil)
	return hashed
}

// 生成新的区块
func (node *Node) GenerateNewBlock(lastBlock Block, data []byte) Block {
	var newBlock = Block{
		Index:     lastBlock.Index + 1,
		PreHash:   lastBlock.Hash,
		Hash:      "",
		TimeStamp: time.Now().String(),
		Data:      data,
		delegate:  nil,
	}
	newBlock.Hash = hex.EncodeToString(blockHash(newBlock))
	newBlock.delegate = node
	return newBlock
}

// 创建10个节点
var NodeAddr = make([]Node, 10)

// 创建节点
func CreateNode() {
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("节点%d 票数", i)
		// 初始化时票数为0
		NodeAddr[i] = Node{
			Name:  name,
			Votes: 0,
		}
	}
}

// 简单模拟投票
func Vote() {
	for i := 0; i < 10; i++ {
		rand.NewSource(time.Now().Unix())
		time.Sleep(100000)
		vote := rand.Intn(10000)
		NodeAddr[i].Votes = vote
		fmt.Printf("节点[%d] 票数[%d]\n", i, vote)
	}
}

// 一共10个节点，选出票数最多的前三个
func SortNodes() []Node {
	n := NodeAddr
	for i := 0; i < len(n); i++ {
		for j := 0; j < len(n)-1; j++ {
			if n[j].Votes < n[j+1].Votes {
				n[j], n[j+1] = n[j+1], n[j]
			}
		}
	}
	return n[:3]
}

func main() {
	// 初始化10个全节点
	CreateNode()
	fmt.Println("创建的节点列表：")
	fmt.Println(NodeAddr)
	fmt.Println("节点票数")
	//投票
	Vote()
	//选出前三
	nodes := SortNodes()
	fmt.Println("获胜者")
	//创世区块
	first := firstBlock()
	lastBlock := first
	fmt.Println("开始生成区块")
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("[%s %d]生成新的区块\n", nodes[i].Name, nodes[i].Votes)
		lastBlock = nodes[i].GenerateNewBlock(lastBlock, []byte(fmt.Sprintf("new Block %d", i)))
	}
}
