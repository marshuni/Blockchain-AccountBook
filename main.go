package main

import (
	"accountbook/pkg/wallet"
	"fmt"
)

func main() {
	myWallet := wallet.NewWallet()
	fmt.Printf("公钥：%x\n", myWallet.PublicKey)
	fmt.Printf("公钥Hash：%x\n", wallet.HashPubKey(myWallet.PublicKey))
	fmt.Println("比特币地址：", myWallet.GetAddress())
}
