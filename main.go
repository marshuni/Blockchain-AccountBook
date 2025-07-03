package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/marshuni/Blockchain-AccountBook/pkg/accountbook"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/wallet"
)

var (
	ab         *accountbook.AccountBook
	walletList []*wallet.Wallet
)

func main() {
	ab = accountbook.NewAccountBook("./database/data.db")
	walletList = []*wallet.Wallet{}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n==== 区块链账本菜单 ====")
		fmt.Println("1. 创建新钱包")
		fmt.Println("2. 查看所有钱包地址")
		fmt.Println("3. 查看钱包余额")
		fmt.Println("4. 转账")
		fmt.Println("5. 添加Coinbase交易")
		fmt.Println("6. 查看所有交易")
		fmt.Println("7. 打印区块链")
		fmt.Println("0. 退出")
		fmt.Print("请选择操作: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			w := ab.NewWallet()
			walletList = append(walletList, w)
			fmt.Println("新钱包已创建，地址：", ab.GetAddress(w))
		case "2":
			if len(walletList) == 0 {
				fmt.Println("暂无钱包，请先创建。")
			} else {
				fmt.Println("所有钱包地址：")
				for i, w := range walletList {
					fmt.Printf("%d: %s\n", i, ab.GetAddress(w))
				}
			}
		case "3":
			fmt.Print("请输入钱包编号或地址: ")
			addr := readWalletAddr(reader)
			balance := ab.GetBalance(addr)
			fmt.Printf("地址 %s 的余额为: %d\n", addr, balance)
		case "4":
			if len(walletList) < 1 {
				fmt.Println("请先创建钱包。")
				continue
			}
			fmt.Print("请输入转出钱包编号: ")
			fromIdx := readWalletIndex(reader)
			if fromIdx < 0 || fromIdx >= len(walletList) {
				fmt.Println("钱包编号无效。")
				continue
			}
			fmt.Print("请输入收款地址: ")
			toAddr := readWalletAddr(reader)
			fmt.Print("请输入转账金额: ")
			amountStr, _ := reader.ReadString('\n')
			amountStr = strings.TrimSpace(amountStr)
			amount, err := strconv.Atoi(amountStr)
			if err != nil || amount <= 0 {
				fmt.Println("金额无效。")
				continue
			}
			newTx, err := ab.CreateTransaction(ab.GetAddress(walletList[fromIdx]), toAddr, amount, walletList[fromIdx])
			if err != nil {
				fmt.Println("转账失败：", err)
				continue
			}
			ab.AddBlock([]*tx.Transaction{newTx}, "")
			fmt.Println("转账交易已打包进新区块，交易ID:", fmt.Sprintf("%x", newTx.ID))
		case "5":
			fmt.Print("请输入接收Coinbase奖励的钱包编号: ")
			idx := readWalletIndex(reader)
			if idx < 0 || idx >= len(walletList) {
				fmt.Println("钱包编号无效。")
				continue
			}
			cbTx := ab.NewCoinbaseTx(ab.GetAddress(walletList[idx]), "")
			ab.AddBlock([]*tx.Transaction{cbTx}, "")
			fmt.Println("Coinbase交易已添加。")
		case "6":
			fmt.Println("区块链所有交易：")
			printAllTransactions()
		case "7":
			ab.PrintChain()
		case "0":
			fmt.Println("退出程序。")
			return
		default:
			fmt.Println("无效操作，请重新选择。")
		}
	}
}

// 辅助函数：读取钱包地址或编号
func readWalletAddr(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	// 如果输入为数字，视为钱包编号
	if idx, err := strconv.Atoi(input); err == nil && idx >= 0 && idx < len(walletList) {
		return ab.GetAddress(walletList[idx])
	}
	return input
}

// 辅助函数：读取钱包编号
func readWalletIndex(reader *bufio.Reader) int {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return idx
}

// 打印所有区块链上的交易
func printAllTransactions() {
	for i, block := range ab.Chain.Blocks {
		fmt.Printf("区块 #%d:\n", i)
		for _, t := range block.Transactions {
			t.PrintDetails()
		}
	}
}
