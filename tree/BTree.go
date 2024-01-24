package tree

import (
	"fmt"
	"math"
)

type BTree struct {
	Value int
	Left  *BTree
	Right *BTree
}

func (t *BTree) String() {
	fmt.Printf("二叉树:值是%d\n", t.Value)
	var leftValue interface{}
	if t.Left != nil {
		leftValue = t.Left.Value
	}
	var rightValue interface{}
	if t.Right != nil {
		rightValue = t.Right.Value
	}
	fmt.Printf("左节点:%v   右节点:%v \n", leftValue, rightValue)
}

// max 以字符串的形式打印
func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

// Level 获取层高
func (t *BTree) Level() int {
	if t == nil {
		return 0
	}
	return max(t.Left.Level(), t.Right.Level()) + 1
}

func printBlanks(count float64) {
	for i := 0.0; i < count; i++ {
		fmt.Print(" ")
	}
}

// PrintBTree 打印树
func PrintBTree(nodes BTrees, maxlevel int, currlevel int) {
	if len(nodes) == 0 || nodes.IsAllNil() {
		return
	}
	floor := maxlevel - currlevel
	// 元素左边的空格数
	leftBlanks := math.Pow(2.0, float64(floor)) - 1
	// 元素左边的空格数
	rightBlanks := math.Pow(2.0, float64(floor+1)) - 1
	// 先打左边空格
	printBlanks(leftBlanks)

	newNodes := make(BTrees, 0)
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
				// 右斜线
				fmt.Print("\\")
			} else {
				printBlanks(1)
			}
			printBlanks(2*lineNums - i)
		}
		fmt.Print("\n")
	}
	PrintBTree(newNodes, maxlevel, currlevel+1)
}

// Preorder 先序遍历
func (t *BTree) Preorder() {
	if t == nil {
		return
	}
	fmt.Printf("%d->", t.Value) //先打印出值
	t.Left.Preorder()           //递归
	t.Right.Preorder()          //递归
}

// Inorder 中序遍历
func (t *BTree) Inorder() {
	if t == nil {
		return
	}
	// 递归
	t.Left.Inorder()
	// 先打印出值
	fmt.Printf("%d->", t.Value)
	// 递归
	t.Right.Inorder() //递归
}

// Postorder 后序遍历
func (t *BTree) Postorder() {
	if t == nil {
		return
	}
	// 递归
	t.Left.Postorder()
	// 递归
	t.Right.Postorder()
	// 先打印出值
	fmt.Printf("%d->", t.Value)
}

func (t *BTree) ConnectLeft(treeOrValue interface{}) *BTree {
	if bt, ok := treeOrValue.(*BTree); ok {
		t.Left = bt
	} else if v, ok := treeOrValue.(int); ok {
		t.Left = NewBTree(v)
	}
	return t
}
func (t *BTree) ConnectRight(treeOrValue interface{}) *BTree {
	if bt, ok := treeOrValue.(*BTree); ok {
		t.Right = bt
	} else if v, ok := treeOrValue.(int); ok {
		t.Right = NewBTree(v)
	}
	return t

}

func NewBTree(value int) *BTree {
	return &BTree{Value: value}
}

// BTrees 二叉树集合类型
type BTrees []*BTree

func (b BTrees) String() {
	for _, bt := range b {
		bt.String()
	}
	fmt.Printf("当前一共有%d个节点", len(b))
}

// IsAllNil 判断子元素是否都为nil
func (b BTrees) IsAllNil() bool {
	for _, bt := range b {
		if bt != nil {
			return false
		}
	}
	return true
}

func NewBTrees(values ...int) BTrees {
	btrees := make(BTrees, len(values))
	for index, v := range values {
		btrees[index] = NewBTree(v)
	}
	return btrees
}
