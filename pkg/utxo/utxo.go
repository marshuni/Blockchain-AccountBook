package utxo

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"slices"

	"github.com/marshuni/Blockchain-AccountBook/pkg/blockchain"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/tx"
	"github.com/marshuni/Blockchain-AccountBook/pkg/core/wallet"
)

// UTXO集
type UTXOSet struct {
	Blockchain *blockchain.Blockchain
}

// 查询未花费输出专用数据结构
type UTXOOutput struct {
	TxID  []byte
	Vout  int
	Value int // 金额
}

// 查找某地址所有未花费输出（查询余额用）
func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []UTXOOutput {
	var utxos []UTXOOutput
	spentTXOs := make(map[string][]int) // 记录已经消费掉的输出，使用Output ID作为键，索引作为值
	// 两层循环遍历每一笔交易
	for _, block := range u.Blockchain.Blocks {
		for _, t := range block.Transactions {
			// 处理输入，标记已花费的输出
			if !t.IsCoinbase() {
				for _, vin := range t.Inputs {
					if bytes.Equal(wallet.HashPubKey(vin.PubKey), pubKeyHash) {
						spentTXOs[string(vin.Txid)] = append(spentTXOs[string(vin.Txid)], vin.Vout)
					}
				}
			}
		}
	}

	for _, block := range u.Blockchain.Blocks {
		for _, t := range block.Transactions {
			txIDStr := string(t.ID)
			// 处理输出，将当前地址下未花费的输出加入Output
			for idx, out := range t.Outputs {
				if bytes.Equal(out.PubKeyHash, pubKeyHash) {
					// 检查该输出是否已被花费
					spentFlag := false
					if spentTXOs[txIDStr] != nil {
						if slices.Contains(spentTXOs[txIDStr], idx) {
							// 若当前Output被消费掉，直接跳过
							spentFlag = true
						}
					}
					if spentFlag {
						continue
					}
					utxos = append(utxos, UTXOOutput{
						TxID:  t.ID,
						Vout:  idx,
						Value: out.Value,
					})
				}
			}
		}
	}
	return utxos

}

// 返回足以覆盖amount的未花费输出
func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, []UTXOOutput) {
	accumulated := 0
	var selectedUTXOs []UTXOOutput
	utxos := u.FindUTXO(pubKeyHash)
	for _, out := range utxos {
		if accumulated >= amount {
			// 选用的Output够用了就停止
			break
		}
		accumulated += out.Value
		selectedUTXOs = append(selectedUTXOs, out)
	}
	return accumulated, selectedUTXOs
}

// 构造新交易
func (u *UTXOSet) CreateTransaction(from, to string, amount int, w *wallet.Wallet) (*tx.Transaction, error) {
	pubKeyHash := wallet.GetPubKeyHashFromAddress(from)
	accumulated, validOutputs := u.FindSpendableOutputs(pubKeyHash, amount)
	if accumulated < amount {
		return nil, errors.New("余额不足")
	}
	var inputs []tx.TXInput
	var outputs []tx.TXOutput

	// 构造输入
	for _, utxo := range validOutputs {
		input := tx.TXInput{
			Txid:      utxo.TxID,
			Vout:      utxo.Vout,
			Signature: nil, // 签名后面再加
			PubKey:    w.PublicKey,
		}
		inputs = append(inputs, input)
	}

	// 构造输出
	outputs = append(outputs, tx.TXOutput{
		Value:      amount,
		PubKeyHash: wallet.GetPubKeyHashFromAddress(to),
	})
	if accumulated > amount {
		// 找零
		outputs = append(outputs, tx.TXOutput{
			Value:      accumulated - amount,
			PubKeyHash: pubKeyHash,
		})
	}

	newTx := &tx.Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	newTx.ID = newTx.CalcID()
	// 签名
	u.SignTransaction(newTx, w.PrivateKey)
	return newTx, nil
}

// 签名交易
func (u *UTXOSet) SignTransaction(t *tx.Transaction, privKey *ecdsa.PrivateKey) {
	if t.IsCoinbase() {
		return
	}
	for idx, vin := range t.Inputs {
		prevTx := u.Blockchain.FindTx(vin.Txid)
		if prevTx == nil {
			continue
		}
		// 只对当前输入引用的输出做签名
		// 签名内容为PubKeyHash+TxID
		dataToSign := append(prevTx.Outputs[vin.Vout].PubKeyHash, t.ID...)
		hash := sha256.Sum256(dataToSign)
		r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
		if err != nil {
			// 打印签名错误
			println("签名失败:", err.Error())
			continue
		}
		signature := append(r.Bytes(), s.Bytes()...)
		t.Inputs[idx].Signature = signature
	}
}
