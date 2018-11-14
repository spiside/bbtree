package bbtree

import "fmt"

const (
	blue  = true
	black = false
)

type Comparable interface {
	More(than Comparable) bool
}

type node struct {
	key         Comparable
	value       interface{}
	left, right *node
	size        int
	parentColor bool
}

func (n node) String() string {
	return fmt.Sprintf("%v", n.key)
}

func (n node) more(than *node) bool {
	return n.key.More(than.key)
}

func size(n *node) int {
	if n == nil {
		return 0
	}
	return n.size
}

func newNode(key Comparable, value interface{}) *node {
	return &node{key, value, nil, nil, 1, blue}
}

func isBlue(n *node) bool {
	if n == nil {
		return false
	}
	return n.parentColor
}
