package tx

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"accountbook/pkg/wallet"
)

// UTXO 结构
type TXInput struct {
	Txid      []byte // 引用的交易ID
	Vout      int    // 引用的Vout索引
	Signature []byte // 用交易发起方的私钥签名，用于确定发起人确实拥有这笔钱
	PubKey    []byte // 交易发起方的原始公钥（并非比特币地址）
}
type TXOutput struct {
	Value      int    // 金额
	PubKeyHash []byte // 交易输出方地址，即交易完成后实际拥有这笔钱的一方
}
type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

// 计算交易ID(Hash)
func (tx *Transaction) CalcID() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	// 使用Gob库而非Json对交易进行序列化，效率更高
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(txCopy)
	if err != nil {
		panic(err)
	}
	hash = sha256.Sum256(buf.Bytes())
	return hash[:]
}

// 创建Coinbase交易
// Coinbase交易由挖矿产生，不涉及到用户主动的交易操作，故不放置到utxo模块
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, []byte{}, []byte(data)}
	txout := TXOutput{100, wallet.GetPubKeyHashFromAddress(to)}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.CalcID()
	return &tx
}

// 判断交易是否为 coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].Txid) == 0 && tx.Inputs[0].Vout == -1
}

// 输出交易细节
func (tx *Transaction) PrintDetails() {
	fmt.Printf("Transaction ID: %x\n", tx.ID)
	if tx.IsCoinbase() {
		fmt.Println("Coinbase Transaction")
		fmt.Printf("%s\n", tx.Inputs[0].PubKey)
	} else {
		fmt.Println("Inputs:")
		for _, in := range tx.Inputs {
			fmt.Printf("  Txid: %x\n", in.Txid)
			fmt.Printf("  Vout: %d\n", in.Vout)
			fmt.Printf("  PubKey: %x\n", in.PubKey)
		}
	}
	fmt.Println("Outputs:")
	for _, out := range tx.Outputs {
		fmt.Printf("  Value: %d\n", out.Value)
		fmt.Printf("  PubKeyHash: %x\n", out.PubKeyHash)
	}
}
