package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// Block
type Block struct {
	Index     int
	TimeStamp string
	BPM       int
	HashCode  string
	PreHash   string
	// 区块验证者
	Validator string
}

// 声明链
var Blockchain []Block

// 临时缓冲区
var tempBlocks []Block

// 声明候选人
// 任何一个节点提议一个新块时，将它发送到这个信道
var candidateBlocks = make(chan Block)

// 公告的信道,用于网络广播的内容
var announcements = make(chan string)

// 锁
var mutex = &sync.Mutex{}

// 验证者列表
// 存储节点地址和tokens
var validators = make(map[string]int)

// 生成区块
func generateBlock(lastBlock Block, BPM int, address string) Block {
	var newBlock Block
	newBlock.Index = lastBlock.Index + 1
	newBlock.TimeStamp = time.Now().String()
	newBlock.PreHash = lastBlock.HashCode
	newBlock.BPM = BPM
	newBlock.Validator = address
	newBlock.HashCode = GenerateHashValue(newBlock)
	return newBlock
}
func GenerateHashValue(block Block) string {
	var hashCode = block.PreHash + block.TimeStamp + block.Validator + strconv.Itoa(block.BPM) + strconv.Itoa(block.Index)

	return calculateHash(hashCode)
}
func calculateHash(s string) string {
	var sha = sha256.New()
	sha.Write([]byte(s))
	hashed := sha.Sum(nil)
	return hex.EncodeToString(hashed)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	genesisBlock := Block{}
	genesisBlock = Block{
		Index:     0,
		TimeStamp: time.Now().String(),
		BPM:       0,
		HashCode:  GenerateHashValue(genesisBlock),
		PreHash:   "",
		Validator: "",
	}
	spew.Dump(genesisBlock)
	//将创世区块，加入到区块数组中
	Blockchain = append(Blockchain, genesisBlock)
	port := os.Getenv("PORT")
	// 启动服务器
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	//打印监听到的端口
	log.Println("HTTP Server Listenin on port:", port)
	//释放资源
	defer server.Close()

	go func() {
		for cadidate := range candidateBlocks {
			// 锁
			mutex.Lock()
			//当候选人中有数据，添加到零时缓冲区
			tempBlocks = append(tempBlocks, cadidate)
			mutex.Unlock()
		}
	}()
	// 查谁去挖矿
	go func() {
		for {
			// 根据tokens的个数去做重划分
			pickWinner()
		}
	}()
	// 接收验证者节点的连接
	for {
		//等待终端连接
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//连上的情况下，处理终端发来的消息
		go handleConn(conn)
	}
}

// Pos 的主要逻辑
// 实现获取记账权的节点
// 根据令牌数
func pickWinner() {
	//休眠30秒
	time.Sleep(30 * time.Second)
	//锁
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	// 声明一个彩票池地址，存放每个验证者地址
	lotteryPool := []string{}
	//如果临时缓冲去又验证者
	if len(temp) > 0 {
		//有验证者
		//根据被标记的令牌的数量对他们进行加权
		// 遍历temp
	OUTER:
		for _, block := range temp {
			// 查看是否已经在彩票池
			for _, node := range lotteryPool {
				if block.Validator == node {
					// 跳出
					continue OUTER
				}
			}

			// 锁
			mutex.Lock()
			//地址和tokens
			setValidators := validators
			mutex.Unlock()

			// 获取验证者的token的个数
			//k：当前验证者的tokens个数
			k, ok := setValidators[block.Validator]
			if ok {
				//向彩票池加入K条数据
				// 将所以得验证者添加到一个数组中
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// 设置随机种子，保证随机性
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		// 通过随机值[0,彩票池长度)随机获取记账权的节点
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]
		// 把获胜者的区块添加到整条区块链上
		//通知其他节点相关胜利者的消息
		for _, block := range temp {
			//判断是否是被选中的孩子
			if block.Validator == lotteryWinner {
				mutex.Lock()
				// 将获胜者的区块添加到区块链中
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				// 广播消息
				for _ = range validators {
					// 将获胜者的地址放到公告里
					announcements <- "\nvalidator:" + lotteryWinner + "\n"
				}
				break
			}
		}
	}

	// 临时缓冲区位空的情况
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}

// 处理终端发来的信息
func handleConn(conn net.Conn) {
	// 释放资源
	defer conn.Close()
	go func() {
		for {
			//打印获胜者的消息
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()
	// 验证者的地址
	var addr string
	// 验证者输入拥有的tokens
	io.WriteString(conn, "Enter token balance:")
	// 接收并处理输入的信息
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		//获取输入的数据，转成int
		//获取余额，持币数量
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number:%v", scanBalance.Text(), err)
			return
		}
		// 生成验证者的地址
		addr = calculateHash(time.Now().String())
		//将验证者的地址和token存到validator
		validators[addr] = balance
		fmt.Println(validators)
		break
	}
	// 输入交易信息
	io.WriteString(conn, "\nEnter a new BPM:")
	scanBPM := bufio.NewScanner(conn)
	go func() {
		// 多次输入交易信息
		for {
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				if err != nil {
					log.Printf("%v bpm not a number:%v", scanBPM.Text(), err)
					//提出改节点
					//从map中移除验证者信息
					// 对恶意节点的惩罚
					delete(validators, addr)
					conn.Close()
				}
				// 取到区块
				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()
				// 创建新的区块
				newBlock := generateBlock(oldLastIndex, bpm, addr)
				// 验证区块
				if isBlockValid(newBlock, oldLastIndex) {
					//验证通过，将新的区块放到通道
					candidateBlocks <- newBlock
				}
			}
		}
	}()
	// 周期性打印区块消息
	for {
		time.Sleep(time.Second * 10)
		mutex.Lock()
		//json 输出
		out, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal("json.Marshal error", err)
			return
		}
		// 输出
		io.WriteString(conn, string(out)+"\n")
	}
}

func isBlockValid(newBlock, oldBlock Block) bool {
	//检查index
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	//preHash
	if oldBlock.HashCode != newBlock.PreHash {
		return false
	}
	//再次验证hash
	if GenerateHashValue(newBlock) != newBlock.HashCode {
		return false
	}
	return true
}
