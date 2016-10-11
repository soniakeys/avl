// Public domain

package avl

import (
	"fmt"
	"testing"
)

func (t *Tree) Dump() {
	t.root.dump("")
}

func (root *node) dump(indent string) {
	fmt.Print(indent, root, "\n")
	if root != nil {
		indent += "  "
		root.link[0].dump(indent)
		root.link[1].dump(indent)
	}
}

type intNode struct {
	intKey
	string
}

type intKey int

func (k intKey) Less(k2 Key) bool    { return k < k2.(intNode).intKey }
func (k intKey) Equal(k2 Key) bool   { return k == k2.(intNode).intKey }
func (k intKey) Greater(k2 Key) bool { return k > k2.(intNode).intKey }

func Test(t *testing.T) {
	var a Tree
	fmt.Println("Empty tree:")
	a.Dump()

	fmt.Println("\nInsert test:")
	a.Insert(intNode{3, "three"})
	a.Insert(intNode{1, "one"})
	a.Insert(intNode{4, "four"})
	a.Insert(intNode{1, "uno"})
	a.Insert(intNode{5, "five"})
	a.Dump()

	fmt.Println("\nFind test:")
	fmt.Println(a.Find(intKey(3)))
	fmt.Println(a.Find(intKey(1)))
	fmt.Println(a.Find(intKey(4)))
	fmt.Println(a.Find(intKey(5)))

	fmt.Println("\nRemove test:")
	a.Remove(intKey(3))
	a.Remove(intKey(1))
	a.Dump()
}
