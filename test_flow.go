// go
package main

import (
	"fmt"

	"github.com/marshuni/Blockchain-AccountBook/pkg/blockchain"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/wallet"
	"github.com/marshuni/Blockchain-AccountBook/pkg/utxo"
)

func TestUTXOFlow() {
	// 1. 初始化区块链和UTXO集
	fmt.Println("【1. 初始化区块链和UTXO集】")
	chain := blockchain.NewBlockchain("data.test.db")
	utxoSet := utxo.UTXOSet{Blockchain: chain}

	// 2. 创建两个钱包A、B
	fmt.Println("【2. 创建两个钱包A、B】")
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	addrA := walletA.GetAddress()
	addrB := walletB.GetAddress()
	fmt.Println("    A地址:", addrA)
	fmt.Println("    B地址:", addrB)

	// 3. 创建Coinbase交易给B，打包进区块链
	fmt.Println("【3. 创建Coinbase交易给B，打包进区块链】")
	coinbaseTx := tx.NewCoinbaseTX(addrB, "Hello")
	pool := blockchain.TxPool{}
	pool.AddTx(coinbaseTx)
	chain.AddBlock(&pool, addrA)
	fmt.Println("    添加Coinbase交易到区块链，A应获得100")

	// 4. 查询A余额
	fmt.Println("【4. 查询A余额】")
	pubKeyHashA := wallet.GetPubKeyHashFromAddress(addrA)
	utxosA, _ := utxoSet.FindSpendableOutputs(pubKeyHashA, 1000)
	fmt.Printf("    A所有UTXO: %+v\n", utxoSet.FindUTXO(pubKeyHashA))
	fmt.Printf("    A累计余额: %d\n", utxosA)
	if utxosA < 100 {
		fmt.Printf("    A余额不足，期望100，实际%d", utxosA)
		return
	}

	// 5. A向B转账40，构造交易，签名，验证签名
	fmt.Println("【5. A向B转账40，构造交易，签名，验证签名】")
	txAB, err := utxoSet.CreateTransaction(addrA, addrB, 40, walletA)
	if err != nil {
		fmt.Printf("    创建A->B交易失败: %v", err)
		return
	}
	fmt.Println("    A->B 40交易创建成功，交易ID:", fmt.Sprintf("%x", txAB.ID))
	if !txAB.VerifyTransaction() {
		fmt.Println("    A->B 交易签名验证失败")
		return
	}
	fmt.Println("    A->B 交易签名验证通过")

	// 6. 打包A->B交易进新区块
	fmt.Println("【6. 打包A->B交易进新区块】")
	pool2 := blockchain.TxPool{}
	pool2.AddTx(txAB)
	chain.AddBlock(&pool2, "")
	fmt.Println("    A->B交易已打包进新区块")

	// 7. 查询A、B余额
	fmt.Println("【7. 查询A、B余额】")
	utxosA2, _ := utxoSet.FindSpendableOutputs(pubKeyHashA, 1000)
	pubKeyHashB := wallet.GetPubKeyHashFromAddress(addrB)
	utxosB, _ := utxoSet.FindSpendableOutputs(pubKeyHashB, 1000)
	fmt.Printf("    A所有UTXO: %+v\n", utxoSet.FindUTXO(pubKeyHashA))
	fmt.Printf("    B所有UTXO: %+v\n", utxoSet.FindUTXO(pubKeyHashB))
	fmt.Printf("    A累计余额: %d\n", utxosA2)
	fmt.Printf("    B累计余额: %d\n", utxosB)
	if utxosA2 != 60 {
		fmt.Printf("    A余额错误，期望60，实际%d", utxosA2)
		return
	}
	if utxosB != 140 {
		fmt.Printf("    B余额错误，期望140，实际%d", utxosB)
		return
	}

	// 8. 验证A->B交易签名
	fmt.Println("【8. 验证A->B交易签名】")
	if !txAB.VerifyTransaction() {
		fmt.Println("    A->B 交易签名验证失败")
	}
	fmt.Println("    A->B 交易签名再次验证通过")
}
