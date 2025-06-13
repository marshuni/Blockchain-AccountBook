package main

import (
	"crypto/sha256"
	"fmt"
)

type BlockData struct {
	content string
}

// Merkle树节点定义
type MerkleNode struct {
	data        *BlockData
	hash        [32]byte
	left_child  *MerkleNode
	right_child *MerkleNode
}

// 从数据块构建Merkle树，返回根
func Create_tree(datas []BlockData) MerkleNode {
	var nodes []MerkleNode
	// 遍历所有数据块，创建叶节点
	for _, data := range datas {
		var newNode MerkleNode

		newNode.data = &data
		update_hash(&newNode)

		nodes = append(nodes, newNode)
	}
	return build_tree(nodes)[0]
}

// 递归逐层构建Merkle树
func build_tree(sons []MerkleNode) []MerkleNode {
	var fathers []MerkleNode
	// 相邻节点配对，创建它们的父节点
	for i := 0; i < len(sons); i += 2 {
		var newNode MerkleNode

		newNode.left_child = &sons[i]
		if i+1 < len(sons) {
			newNode.right_child = &sons[i+1]
		}
		update_hash(&newNode)

		fathers = append(fathers, newNode)
	}

	if len(fathers) == 1 {
		return fathers
	} else {
		return build_tree(fathers)
	}
}

// 递归打印各节点的哈希值
func Print_tree(now MerkleNode, layer int) {
	for range layer {
		fmt.Print("  ")
	}
	if now.data != nil {
		fmt.Printf("Data:\"%s\":%x\n", now.data.content, now.hash)
		return
	}

	fmt.Printf("Node:%x\n", now.hash)
	if now.left_child != nil {
		Print_tree(*now.left_child, layer+1)
	}
	if now.right_child != nil {
		Print_tree(*now.right_child, layer+1)
	}
}

// 根据子节点信息，计算当前节点Hash
func update_hash(node *MerkleNode) {
	hash := sha256.New()

	// 存在data，说明是叶节点
	// 否则根据左右子节点的哈希计算
	if node.data != nil {
		// 若指向数据块，根据content属性计算哈希值
		hash.Write([]byte(node.data.content))
	} else {
		hash.Write(node.left_child.hash[:])
		if node.right_child != nil {
			hash.Write(node.right_child.hash[:])
		}
	}
	copy(node.hash[:], hash.Sum(nil))
}

// func main() {
// 	datas := [5]DataBlock{{"ha"}, {"hello"}, {"world"}, {"aaa"}, {"hello~"}}
// 	root := create_tree(datas[:])
// 	fmt.Printf("Root Hash of Merkle Tree:%x\n", root.hash)

// 	print_tree(root, 0)
// }

/*
参考链接：
	- https://blog.csdn.net/qq756684177/article/details/81518823
	- https://www.cnblogs.com/X-knight/p/9142622.html
	- https://zhuanlan.zhihu.com/p/666478154
	- https://www.cnblogs.com/wanghui-garcia/p/10452431.html
*/
