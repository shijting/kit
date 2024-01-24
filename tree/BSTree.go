package tree

import (
	"fmt"
	"math"
)

type BSTree struct {
	Value int
	Left  *BSTree
	Right *BSTree
}

type BSTrees []*BSTree

// IsAllNil 判断子元素是否都为nil
func (t BSTrees) IsAllNil() bool {
	for _, bt := range t {
		if bt != nil {
			return false
		}
	}
	return true
}

func NewBSTree(value int) *BSTree {
	return &BSTree{Value: value}
}

func (t *BSTree) Level() int {
	if t == nil {
		return 0
	}
	return max(t.Left.Level(), t.Right.Level()) + 1
}

// PrintBSTree 打印树
func PrintBSTree(nodes BSTrees, maxlevel int, currlevel int) {
	if len(nodes) == 0 || nodes.IsAllNil() {
		return
	}

	floor := maxlevel - currlevel

	// 元素左边的空格数
	leftBlanks := math.Pow(2.0, float64(floor)) - 1

	//元素左边的空格数
	rightBlanks := math.Pow(2.0, float64(floor+1)) - 1

	// 先打左边空格
	printBlanks(leftBlanks)

	newNodes := make(BSTrees, 0)
	for _, node := range nodes {
		if node != nil {
			fmt.Print(node.Value)
			newNodes = append(newNodes, node.Left, node.Right)
		} else {
			// 打印一个空格
			printBlanks(1)
			newNodes = append(newNodes, nil, nil)
		}
		// 打右边空格
		printBlanks(rightBlanks)
	}
	fmt.Print("\n")

	//画线
	lineNums := math.Pow(2, float64(floor-1))
	for i := 1.0; i <= lineNums; i++ {
		for _, node := range nodes {
			// 左边线做空格
			printBlanks(leftBlanks - i)

			if node == nil {
				printBlanks(lineNums*2 + i + 1)
				continue
			}
			if node.Left != nil {
				fmt.Print("/")
			} else {
				printBlanks(1)
			}

			// 左边线的右空格
			printBlanks(2*i - 1)
			if node.Right != nil {
				// 右边线
				fmt.Print("\\")
			} else {
				printBlanks(1)
			}

			printBlanks(2*lineNums - i)
		}
		fmt.Print("\n")
	}
	PrintBSTree(newNodes, maxlevel, currlevel+1)
}

// AddBSTree 递归加入节点
func AddBSTree(bst *BSTree, node *BSTree) *BSTree {
	if node == nil { //没有子节点了
		return bst
	}

	if bst.Value < node.Value {
		// 左边节点都小于根
		node.Left = AddBSTree(bst, node.Left)
	} else if bst.Value > node.Value {
		// 右边节点都大于根
		node.Right = AddBSTree(bst, node.Right)
	}

	return node
}

func SearchNodeWithParent(value int, node *BSTree, parent ...interface{}) (*BSTree, *BSTree, string) {
	if node == nil {
		return nil, nil, ""
	}

	if value < node.Value {
		return SearchNodeWithParent(value, node.Left, node, "left")
	} else if value > node.Value {
		return SearchNodeWithParent(value, node.Right, node, "right")
	} else {
		if len(parent) == 0 { //没有传参数
			return node, nil, ""
		}
		return node, parent[0].(*BSTree), parent[1].(string)

	}
}

// SearchNode 搜索节点
func SearchNode(value int, node *BSTree) *BSTree {
	if node == nil {
		return nil
	}
	if value < node.Value {
		return SearchNode(value, node.Left)
	} else if value > node.Value {
		return SearchNode(value, node.Right)
	} else {
		return node
	}
}

// ParentNode 获取父节点
func ParentNode(node *BSTree, root *BSTree) *BSTree {
	if node == nil || root == nil || node == root {
		return nil
	}
	// 第二步
	if node == root.Left || node == root.Right {
		return root
	}
	// 第三步 (递归，左递归->右递归
	getLeft := ParentNode(node, root.Left)
	if getLeft != nil {
		return getLeft
	}
	return ParentNode(node, root.Right)

}

// isLeaf 判断是否有子节点
func (t *BSTree) isLeaf() bool {
	if t.Left == nil && t.Right == nil {
		return true
	}
	return false
}

// getSingleChild 返回 当前节点的子节点（要么左 要么右)
func (t *BSTree) getSingleChild() *BSTree {
	if t.Left != nil && t.Right == nil {
		return t.Left
	}
	if t.Left == nil && t.Right != nil {
		return t.Right
	}
	return nil
}

func DeleteNode(value int, treeRoot *BSTree) {
	node, parent, pos := SearchNodeWithParent(value, treeRoot)
	if node == nil { //没有节点
		return
	}

	// 没有子节点，直接删
	if node.isLeaf() {
		if parent == nil {
			*treeRoot = *(*BSTree)(nil)
			return
		}
		if pos == "left" {
			parent.Left = nil
		} else {
			parent.Right = nil
		}
	} else if single := node.getSingleChild(); single != nil { //代表 需要删除的节点只有一个子节点 要么左 要么右
		if parent == nil {
			if node.Left != nil {
				*treeRoot = *node.Left
			} else {
				*treeRoot = *node.Right
			}
			return
		}
		if pos == "left" {
			parent.Left = single
		} else {
			parent.Right = single
		}
	} else {
		successor := MinNode(node.Right)
		DeleteNode(successor.Value, treeRoot)
		node.Value = successor.Value
	}

}

// Preorder 前序
func (t *BSTree) Preorder(result *[]int) {
	if t == nil {
		return
	}
	*result = append(*result, t.Value)
	// 递归
	t.Left.Preorder(result)
	// 递归
	t.Right.Preorder(result)
}

// Inorder 中序
func (t *BSTree) Inorder(result *[]int) {
	if t == nil {
		return
	}

	// 递归
	t.Left.Inorder(result)
	*result = append(*result, t.Value)
	// 递归
	t.Right.Inorder(result)
}

// Postorder 后序
func (t *BSTree) Postorder(result *[]int) {
	if t == nil {
		return
	}
	t.Left.Postorder(result)
	t.Right.Postorder(result)
	*result = append(*result, t.Value)
}

// MaxNode 最大节点
func MaxNode(node *BSTree) *BSTree {
	if node.Right != nil {
		return MaxNode(node.Right)
	} else {
		return node
	}
}

// MinNode 最小节点
func MinNode(node *BSTree) *BSTree {
	if node.Left != nil {
		return MinNode(node.Left)
	} else {
		return node
	}
}
