package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"accountbook/pkg/core/tx"
)

// Merkle树节点定义
type MerkleNode struct {
	Data       *tx.Transaction
	Hash       [32]byte
	LeftChild  *MerkleNode
	RightChild *MerkleNode
}

// 从数据块构建Merkle树，返回根
func CreateTree(datas []*tx.Transaction) MerkleNode {
	var nodes []MerkleNode
	// 遍历所有数据块，创建叶节点
	for _, data := range datas {
		var newNode MerkleNode

		newNode.Data = data
		updateHash(&newNode)

		nodes = append(nodes, newNode)
	}
	return buildTree(nodes)[0]
}

// 递归逐层构建Merkle树
func buildTree(sons []MerkleNode) []MerkleNode {
	var fathers []MerkleNode
	// 相邻节点配对，创建它们的父节点
	for i := 0; i < len(sons); i += 2 {
		var newNode MerkleNode

		newNode.LeftChild = &sons[i]
		if i+1 < len(sons) {
			newNode.RightChild = &sons[i+1]
		}
		updateHash(&newNode)

		fathers = append(fathers, newNode)
	}

	if len(fathers) == 1 {
		return fathers
	} else {
		return buildTree(fathers)
	}
}

// 递归打印各节点的哈希值
func PrintTree(now MerkleNode, layer int) {
	for range layer {
		fmt.Print("  ")
	}
	if now.Data != nil {
		fmt.Printf("交易ID:\"%x\"\n    哈希值:%x\n", now.Data.ID, now.Hash)
		return
	}

	fmt.Printf("Node:%x\n", now.Hash)
	if now.LeftChild != nil {
		PrintTree(*now.LeftChild, layer+1)
	}
	if now.RightChild != nil {
		PrintTree(*now.RightChild, layer+1)
	}
}

// 根据子节点信息，计算当前节点Hash
func updateHash(node *MerkleNode) {
	hash := sha256.New()

	// 存在data，说明是叶节点
	// 否则根据左右子节点的哈希计算
	if node.Data != nil {
		// 若指向数据块，使用gob序列化并计算哈希值
		var serializedData []byte
		buffer := bytes.NewBuffer(serializedData)
		encoder := gob.NewEncoder(buffer)
		err := encoder.Encode(node.Data)
		if err != nil {
			panic(fmt.Sprintf("Failed to serialize data: %v", err))
		}
		hash.Write(buffer.Bytes())
	} else {
		hash.Write(node.LeftChild.Hash[:])
		if node.RightChild != nil {
			hash.Write(node.RightChild.Hash[:])
		}
	}
	copy(node.Hash[:], hash.Sum(nil))
}

/*
参考链接：
	- https://blog.csdn.net/qq756684177/article/details/81518823
	- https://www.cnblogs.com/X-knight/p/9142622.html
	- https://zhuanlan.zhihu.com/p/666478154
	- https://www.cnblogs.com/wanghui-garcia/p/10452431.html
*/
