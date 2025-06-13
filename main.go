package main

import (
	"accountbook/pkg/tx"
	"accountbook/pkg/wallet"

	"accountbook/pkg/merkle"
	"fmt"
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
	merkle.CreateTree([]tx.Transaction{*myCoinbase, *myCoinbase})
}
