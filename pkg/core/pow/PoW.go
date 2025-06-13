package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"time"

	"accountbook/pkg/core/merkle"
	"accountbook/pkg/core/tx"
)

type Block struct {
	Version      uint32
	PreviousHash [32]byte
	MerkleRoot   [32]byte

	Timestamp uint32
	Bits      [4]byte
	Nounce    uint32

	Transactions []*tx.Transaction
}

func NewBlock(previousHash [32]byte, transactions []*tx.Transaction, bits ...[4]byte) Block {
	var newBlock Block

	newBlock.Version = 2
	newBlock.PreviousHash = previousHash
	merkleRoot := merkle.CreateTree(transactions)
	newBlock.MerkleRoot = merkleRoot.Hash

	newBlock.Timestamp = uint32(time.Now().Unix())
	if len(bits) > 0 {
		newBlock.Bits = bits[0]
	} else {
		newBlock.Bits = [4]byte{0x1f, 0x00, 0xff, 0xff}
	}
	newBlock.Nounce = 0

	newBlock.Transactions = transactions

	return newBlock
}

func (block *Block) CalculateHash() [32]byte {
	hash := sha256.New()

	// 数值型和字节数组采用不同的方式写入
	binary.Write(hash, binary.BigEndian, block.Version)
	hash.Write(block.PreviousHash[:])
	hash.Write(block.MerkleRoot[:])
	binary.Write(hash, binary.BigEndian, block.Timestamp)
	hash.Write(block.Bits[:])
	binary.Write(hash, binary.BigEndian, block.Nounce)

	return [32]byte(hash.Sum(nil))
}

// 将bits转换为难度目标值Target
func BitsToTarget(bits [4]byte) [32]byte {
	// 从bits提取系数和指数
	exponent := bits[0]
	coefficient := binary.BigEndian.Uint32(bits[:]) & 0x00ffffff
	// 根据公式计算目标值
	target := new(big.Int).SetUint64(uint64(coefficient))
	target.Lsh(target, 8*(uint(exponent)-3))
	// 将目标值转为 [32]byte
	var targetBytes [32]byte
	copy(targetBytes[32-len(target.Bytes()):], target.Bytes())
	return targetBytes
}

func (block *Block) MineBlock() [32]byte {
	target := BitsToTarget(block.Bits)

	// bytes.Compare(target1, target2)用于比较字典序
	// 若 target1<target2 返回-1，否则返回1
	for hash := block.CalculateHash(); bytes.Compare(hash[:], target[:]) > 0; {
		block.Nounce++
		block.Timestamp = uint32(time.Now().Unix())

		hash = block.CalculateHash()
		// if block.nounce%1000 == 0 {
		// 	fmt.Printf("Current hash: %x\n", hash)
		// }
	}
	return block.CalculateHash()
}
