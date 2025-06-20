package main

import (
	"fmt"

	"github.com/marshuni/Blockchain-AccountBook/pkg/blockchain"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/merkle"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/pow"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/wallet"
)

func main() {
	// 验证钱包可用性
	myWallet := wallet.NewWallet()
	fmt.Printf("公钥：%x\n", myWallet.PublicKey)
	fmt.Printf("公钥Hash：%x\n", wallet.HashPubKey(myWallet.PublicKey))
	myAddress := myWallet.GetAddress()
	fmt.Println("比特币地址：", myAddress)

	// 验证交易模块可用性
	myCoinbase := tx.NewCoinbaseTX(myAddress, "")
	fmt.Println("---------\n创建一个Coinbase交易：")
	myCoinbase.PrintDetails()

	// 验证Merkle树可用性
	myMerkleTreeRoot := merkle.CreateTree([]*tx.Transaction{myCoinbase, myCoinbase})
	fmt.Println("---------\n创建一个Merkle树：")
	merkle.PrintTree(myMerkleTreeRoot, 0)

	// 验证PoW可用性
	var previousHash = [32]byte{
		0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b,
		0x9c, 0xad, 0xbe, 0xcf, 0xd0, 0xe1, 0xf2, 0x03,
		0x14, 0x25, 0x36, 0x47, 0x58, 0x69, 0x7a, 0x8b,
		0x9c, 0xad, 0xbe, 0xcf, 0xd0, 0xe1, 0xf2, 0x33,
	}
	myBlock := pow.NewBlock(previousHash, []*tx.Transaction{myCoinbase, myCoinbase})
	fmt.Println("---------\n打包区块并挖矿：")
	fmt.Printf("当前难度值: %x\n", pow.BitsToTarget(myBlock.Bits))
	fmt.Printf("Block mined: %x\n", myBlock.MineBlock())

	// 区块链
	myChain := blockchain.NewBlockchain("./database/data.db")
	myPool := blockchain.TxPool{}

	myPool.AddTx(myCoinbase)
	myChain.AddBlock(&myPool, myWallet.GetAddress())
	fmt.Println("---------\n区块链测试：")
	myChain.Print()
}
