package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Block 定义区块
type Block struct {
	Index     int
	PreHash   string
	Hash      string
	TimeStamp string
	BPM       int
	// 区块验证者
	Validator string
}

// 定义区块链
var BlockChain []Block

// 生成区块
func GenerateNextBlock(oldBlock Block, BPM int, addr string) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.PreHash = oldBlock.Hash
	newBlock.BPM = BPM
	newBlock.TimeStamp = time.Now().String()
	newBlock.Validator = addr
	//计算hash
	newBlock.Hash = GenerateHashValue(newBlock)
	return newBlock
}

func GenerateHashValue(block Block) string {
	var hashCode = block.PreHash + block.TimeStamp + block.Validator + strconv.Itoa(block.BPM) + strconv.Itoa(block.Index)
	var sha = sha256.New()
	sha.Write([]byte(hashCode))
	hashed := sha.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 定义4个委托人
var delegate = []string{"aaa", "bbb", "ccc", "ddd"}

// 模拟对委托人位置进行随机处理
// 后面让4个委托人轮询挖矿
// 随机位置i处理，被攻击的概率变小，降低风险
func RandDelegate() {
	//rand.Seed()
	rand.NewSource(time.Now().Unix())
	var r = rand.Intn(3)
	t := delegate[r]
	delegate[r] = delegate[3]
	delegate[3] = t
}

func main() {
	fmt.Println(delegate)
	RandDelegate()
	fmt.Println(delegate)
	// 创建创世区块
	var firstBlock Block
	// 将创世区块加入区块链
	BlockChain = append(BlockChain, firstBlock)

	var n = 0
	for {
		// 每30s产生一个区块
		time.Sleep(time.Second * 3)
		var nextBlock = GenerateNextBlock(firstBlock, 1, delegate[n])
		n++
		n = n % len(delegate)
		firstBlock = nextBlock
		BlockChain = append(BlockChain, nextBlock)
		fmt.Println(nextBlock)
	}
}
