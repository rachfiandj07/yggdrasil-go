package utils

import (
	"fmt"
)

type Tree struct {
	val      NodeValue
	parent   *Tree
	children []*Tree
}

func (t *Tree) Val() NodeValue {
	return t.val
}

func (t *Tree) SetVal(n NodeValue) {
	t.val = n
}

func (t *Tree) Parent() (p *Tree, ok bool) {
	if t.parent == nil {
		return t, false
	}
	return t.parent, true
}

func (t *Tree) Children() []*Tree {
	return t.children
}

func (t *Tree) Child(i int) (child *Tree, err error) {
	if i < 0 || i >= len(t.children) {
		return nil, fmt.Errorf("there is no child with index %d", i)
	}
	return t.children[i], nil
}

func (t *Tree) AddChild(n NodeValue) (tChild *Tree) {
	tChild = &Tree{val: n, parent: t}
	t.children = append(t.children, tChild)
	return
}

func NewTree(val NodeValue) *Tree {
	return &Tree{val: val}
}

func (t *Tree) Root() (root *Tree) {
	for root = t; root.parent != nil; root = root.parent {
	}
	return root
}

func (t *Tree) String() string {
	return pencil(t.Root()).String()
}
