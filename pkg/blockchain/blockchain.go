package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"accountbook/pkg/core/pow"
	"accountbook/pkg/core/tx"
)

// 区块链
type Blockchain struct {
	Blocks []*pow.Block
}

// 交易池
type TxPool struct {
	Transactions []*tx.Transaction
}

// 初始化区块链，含创建创世块
func NewBlockchain() *Blockchain {
	genesisBlock := &pow.Block{
		Version:      2,
		PreviousHash: [32]byte{},
		MerkleRoot:   [32]byte{},
		Timestamp:    uint32(time.Now().Unix()),
		Bits:         [4]byte{0x1f, 0x00, 0xff, 0xff},
		Nounce:       0,
	}
	return &Blockchain{
		Blocks: []*pow.Block{genesisBlock},
	}
}

// 打包交易池中的所有区块，挖掘新的区块并添加到链上（自行添加一个Coinbase）
func (bc *Blockchain) AddBlock(p *TxPool, minerAddress string) {
	// 获取交易池中的所有交易
	transactions := p.PopTx()
	if transactions == nil {
		return // 如果没有交易，则不创建新的区块
	}
	if minerAddress != "" {
		// 添加Coinbase块
		coinbaseTx := tx.NewCoinbaseTX(minerAddress, "")
		transactions = append([]*tx.Transaction{coinbaseTx}, transactions...)
	}

	// 获取前一个区块的哈希值
	var previousHash [32]byte
	if len(bc.Blocks) > 0 {
		previousHash = bc.Blocks[len(bc.Blocks)-1].CalculateHash()
	}

	// 使用pow.NewBlock()方法创建新的区块
	// TODO: 难度值调节
	newBlock := pow.NewBlock(previousHash, transactions)

	// 挖掘区块（工作量证明）
	newBlock.MineBlock()

	// 将新挖掘的区块添加到区块链
	bc.Blocks = append(bc.Blocks, &newBlock)
}

// 区块链迭代器
type BlockchainIterator struct {
	currentIndex int
	blockchain   *Blockchain
}

// 获取下一个区块
func (it *BlockchainIterator) Next() (*pow.Block, error) {
	if it.currentIndex < 0 {
		return nil, errors.New("no more blocks")
	}
	block := it.blockchain.Blocks[it.currentIndex]
	it.currentIndex--
	return block, nil
}

// 迭代器，可获取并遍历链上所有交易
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentIndex: len(bc.Blocks) - 1,
		blockchain:   bc,
	}
}

// 寻找特定ID的交易
func (bc *Blockchain) FindTx(TxID []byte) *tx.Transaction {
	for _, block := range bc.Blocks {
		for _, t := range block.Transactions {
			if bytes.Equal(t.ID, TxID) {
				return t
			}
		}
	}
	return nil
}

// 添加新的交易到交易池
func (p *TxPool) AddTx(t *tx.Transaction) bool {
	for _, tx := range p.Transactions {
		if bytes.Equal(tx.ID, t.ID) {
			return false
		}
	}
	p.Transactions = append(p.Transactions, t)
	return true
}

// 返回交易池的所有交易，清空交易池并返回
func (p *TxPool) PopTx() []*tx.Transaction {
	if len(p.Transactions) == 0 {
		return nil
	}
	result := p.Transactions
	p.Transactions = []*tx.Transaction{}
	return result
}

// 打印区块链所有区块及其交易信息
func (bc *Blockchain) Print() {
	for i, block := range bc.Blocks {
		fmt.Printf("Block #%d:\n", i)
		fmt.Printf("  Version: %d\n", block.Version)
		fmt.Printf("  PreviousHash: %x\n", block.PreviousHash)
		fmt.Printf("  MerkleRoot: %x\n", block.MerkleRoot)
		fmt.Printf("  Timestamp: %d\n", block.Timestamp)
		fmt.Printf("  Bits: %x\n", block.Bits)
		fmt.Printf("  Nounce: %d\n", block.Nounce)
		fmt.Printf("  Hash: %x\n", block.CalculateHash())
		fmt.Printf("  Transactions:\n")
		for j, tx := range block.Transactions {
			fmt.Printf("    Transaction #%d:\n", j)
			fmt.Printf("      ID: %x\n", tx.ID)
			fmt.Printf("      Vin:\n")
			for k, vin := range tx.Inputs {
				fmt.Printf("        Vin #%d:\n", k)
				fmt.Printf("          Txid: %x\n", vin.Txid)
				fmt.Printf("          Vout: %d\n", vin.Vout)
				fmt.Printf("          ScriptSig: %s\n", vin.Signature)
			}
			fmt.Printf("      Vout:\n")
			for k, vout := range tx.Outputs {
				fmt.Printf("        Vout #%d:\n", k)
				fmt.Printf("          Value: %d\n", vout.Value)
				fmt.Printf("          ScriptPubKeyHash: %x\n", vout.PubKeyHash)
			}
		}
		fmt.Println()
	}

}

// Gob序列化与反序列化
// func (bc *Blockchain) Serialize() ([]byte, error) {
// 	var result bytes.Buffer
// 	encoder := gob.NewEncoder(&result)
// 	err := encoder.Encode(bc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result.Bytes(), nil
// }

// func DeserializeBlockchain(data []byte) (*Blockchain, error) {
// 	var bc Blockchain
// 	decoder := gob.NewDecoder(bytes.NewReader(data))
// 	err := decoder.Decode(&bc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &bc, nil
// }
