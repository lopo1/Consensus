package main

import (
	"Consensus/pow/block"
	"Consensus/pow/blockchain"
)

func main() {
	var first = block.GenerateFirstBlock("创世区块")
	var second = block.GenerateNextBlock("第二个区块", first)
	var header = blockchain.CreateHeaderNode(&first)
	blockchain.AddNode(&second, header)
	blockchain.ShowNodes(header)
}
