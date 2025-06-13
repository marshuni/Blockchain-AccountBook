package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"
)

type Block struct {
	version      uint32
	previousHash [32]byte
	merkleRoot   [32]byte

	timestamp uint32
	bits      [4]byte
	nounce    uint32
}

func NewBlock(previousHash [32]byte, merkleRoot [32]byte) Block {
	var newBlock Block

	newBlock.version = 2
	newBlock.previousHash = previousHash
	newBlock.merkleRoot = merkleRoot
	newBlock.timestamp = uint32(time.Now().Unix())
	newBlock.bits = [4]byte{0x1f, 0x00, 0xff, 0xff}
	newBlock.nounce = 0

	return newBlock
}

func CalculateHash(block Block) [32]byte {
	hash := sha256.New()

	// 数值型和字节数组采用不同的方式写入
	binary.Write(hash, binary.BigEndian, block.version)
	hash.Write(block.previousHash[:])
	hash.Write(block.merkleRoot[:])
	binary.Write(hash, binary.BigEndian, block.timestamp)
	hash.Write(block.bits[:])
	binary.Write(hash, binary.BigEndian, block.nounce)

	return [32]byte(hash.Sum(nil))
}

// 将bits转换为难度目标值Target
func bitsToTarget(bits [4]byte) [32]byte {
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

func MineBlock(block *Block) {
	target := bitsToTarget(block.bits)

	// bytes.Compare(target1, target2)用于比较字典序
	// 若 target1<target2 返回-1，否则返回1
	for hash := CalculateHash(*block); bytes.Compare(hash[:], target[:]) > 0; {
		block.nounce++
		block.timestamp = uint32(time.Now().Unix())

		hash = CalculateHash(*block)
		// if block.nounce%1000 == 0 {
		// 	fmt.Printf("Current hash: %x\n", hash)
		// }
	}
	fmt.Printf("Block mined: %x\n", CalculateHash(*block))
}

var (
	previousHash = [32]byte{
		0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b,
		0x9c, 0xad, 0xbe, 0xcf, 0xd0, 0xe1, 0xf2, 0x03,
		0x14, 0x25, 0x36, 0x47, 0x58, 0x69, 0x7a, 0x8b,
		0x9c, 0xad, 0xbe, 0xcf, 0xd0, 0xe1, 0xf2, 0x33,
	}
	merkleRoot = [32]byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
		0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01,
	}
)

// func main() {
// 	block := NewBlock(previousHash, merkleRoot)

// 	target := bitsToTarget(block.bits)
// 	fmt.Printf("Target: %x\n", target)
// 	fmt.Printf("Current Hash: %x\n", CalculateHash(block))
// 	MineBlock(&block)
// }
