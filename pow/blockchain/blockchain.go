package blockchain

import (
	"Consensus/pow/block"
	"fmt"
)

// Node 用链表实现区块链的链
// 定义区块链
type Node struct {
	// 指针域
	NextNode *Node
	// 数据域
	Data *block.Block
}

// 创建头节点，保存创世区块
func CreateHeaderNode(data *block.Block) *Node {
	// 先去初始化
	var headerNode = new(Node)
	// 指针域指向nil
	headerNode.NextNode = nil
	// 数据域
	headerNode.Data = data
	//返回头节点，然后再添加节点
	return headerNode
}

// 当挖到矿，添加区块，添加节点
func AddNode(data *block.Block, node *Node) *Node {
	// 创建新的节点
	var newNode = new(Node)
	newNode.Data = data
	newNode.NextNode = nil
	// 连接链表
	node.NextNode = newNode
	return newNode
}

// 查看链表
func ShowNodes(node *Node) {
	n := node
	for {
		if n.NextNode == nil {
			fmt.Println(n.Data)
			break
		} else {
			fmt.Println(n.Data)
			n = n.NextNode
		}
	}
}
