package main

import (
	"fmt"
)

type Node struct {
	Data  int
	Left  *Node
	Right *Node
}

func main() {
	tree := setTree(1)
	tree.Left = setTree(2)
	tree.Right = setTree(3)
	tree.Left.Left = setTree(4)
	tree.Left.Right = setTree(5)

	inOrderTraversal(tree)
}

func setTree(value int) *Node {
	var node = new(Node)
	node.Data = value
	node.Left = nil
	node.Right = nil
	return node
}

func inOrderTraversal(root *Node) {
	if root == nil {
		return
	}
	inOrderTraversal(root.Left)
	fmt.Print(root.Data, "\t")
	inOrderTraversal(root.Right)
}
