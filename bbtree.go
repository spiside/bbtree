package bbtree

import (
	"fmt"
	"strings"
	"sync"
)

type BBTree struct {
	root *node
	mu   sync.Mutex
}

// Node type holds key and value variables for client use.
type Node struct {
	Key   Comparable
	Value interface{}
}

func (n Node) String() string {
	return fmt.Sprintf("{Key: %v, Value: %v}", n.Key, n.Value)
}

// IsEmpty returns a bool indictating whether there are elements
// in the tree.
func (bb *BBTree) IsEmpty() bool {
	return bb.root == nil
}

func (bb *BBTree) rotateRight(currentNode *node) *node {
	// Assert is blue on left.
	newRoot := currentNode.left
	currentNode.left = newRoot.right
	newRoot.right = currentNode

	newRoot.size = currentNode.size
	currentNode.size = size(currentNode.left) + size(currentNode.right) + 1

	newRoot.parentColor = currentNode.parentColor
	currentNode.parentColor = blue

	return newRoot
}

func (bb *BBTree) rotateLeft(currentNode *node) *node {
	// Assert is blue on right.
	newRoot := currentNode.right
	currentNode.right = newRoot.left
	newRoot.left = currentNode

	newRoot.size = currentNode.size
	currentNode.size = size(currentNode.left) + size(currentNode.right) + 1

	newRoot.parentColor = currentNode.parentColor
	currentNode.parentColor = blue
	return newRoot
}

func (bb *BBTree) flipColors(currentNode *node) {
	currentNode.parentColor = !currentNode.parentColor
	currentNode.left.parentColor = !currentNode.left.parentColor
	currentNode.right.parentColor = !currentNode.right.parentColor
}

// fixUp combines the rotate and flip methods to ensure invariants
// during Put() and Remove().
func (bb *BBTree) fixUp(currentNode *node) *node {
	if isBlue(currentNode.left) && !isBlue(currentNode.right) {
		currentNode = bb.rotateRight(currentNode)
	}
	if isBlue(currentNode.right) && isBlue(currentNode.right.right) {
		currentNode = bb.rotateLeft(currentNode)
	}

	if isBlue(currentNode.right) && isBlue(currentNode.left) {
		bb.flipColors(currentNode)
	}

	currentNode.size = size(currentNode.left) + size(currentNode.right) + 1
	return currentNode
}

func (bb *BBTree) moveBlueRight(currentNode *node) *node {
	bb.flipColors(currentNode)

	if isBlue(currentNode.left.right) {
		currentNode.left = bb.rotateLeft(currentNode.left)
		currentNode = bb.rotateRight(currentNode)
		bb.flipColors(currentNode)
	}
	return currentNode
}

func (bb *BBTree) moveBlueLeft(currentNode *node) *node {
	bb.flipColors(currentNode)

	if isBlue(currentNode.right.right) {
		currentNode = bb.rotateLeft(currentNode)
		bb.flipColors(currentNode)
	}
	return currentNode
}

// Put inserts or replaces a given key with an associated value
// in the tree.
func (bb *BBTree) Put(key Comparable, value interface{}) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if value == nil {
		panic("value must not be a nil value")
	}

	putNode := newNode(key, value)
	bb.root = bb.put(bb.root, putNode)
	bb.root.parentColor = black
}

func (bb *BBTree) put(currentNode, putNode *node) *node {
	if currentNode == nil {
		return putNode
	}

	if putNode.more(currentNode) {
		currentNode.right = bb.put(currentNode.right, putNode)
	} else if currentNode.more(putNode) {
		currentNode.left = bb.put(currentNode.left, putNode)
	} else {
		currentNode.value = putNode.value
	}

	return bb.fixUp(currentNode)
}

// Max returns the max Node in the tree.
func (bb *BBTree) Max() *Node {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if bb.IsEmpty() {
		return nil
	}

	return bb.max(bb.root)
}

func (bb *BBTree) max(currentNode *node) *Node {
	if currentNode.right == nil {
		return &Node{currentNode.key, currentNode.value}
	}
	return bb.max(currentNode.right)
}

// Min returns the min Node in the tree.
func (bb *BBTree) Min() *Node {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if bb.IsEmpty() {
		return nil
	}

	return bb.min(bb.root)
}

func (bb *BBTree) min(currentNode *node) *Node {
	if currentNode.left == nil {
		return &Node{currentNode.key, currentNode.value}
	}

	return bb.min(currentNode.left)
}

// RemoveMax removes the max node in the tree.
func (bb *BBTree) RemoveMax() {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if bb.IsEmpty() {
		return
	}

	bb.root = bb.removeMax(bb.root)
	bb.root.parentColor = black
}

func (bb *BBTree) removeMax(currentNode *node) *node {
	if currentNode.right == nil {
		return nil
	}

	if !isBlue(currentNode.right) && !isBlue(currentNode.right.right) {
		currentNode = bb.moveBlueRight(currentNode)
	}
	currentNode.right = bb.removeMax(currentNode.right)

	return bb.fixUp(currentNode)
}

// RemoveMin removes the min node from the tree.
func (bb *BBTree) RemoveMin() {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if bb.IsEmpty() {
		return
	}
	bb.root = bb.removeMin(bb.root)
	bb.root.parentColor = black
}

func (bb *BBTree) removeMin(currentNode *node) *node {
	if isBlue(currentNode.right) {
		currentNode = bb.rotateLeft(currentNode)
	}

	if currentNode.left == nil {
		return nil
	}

	if !isBlue(currentNode.left) && !isBlue(currentNode.left.right) {
		currentNode = bb.moveBlueLeft(currentNode)
	}
	currentNode.left = bb.removeMin(currentNode.left)

	return bb.fixUp(currentNode)
}

// Remove removes a given Comparable key from the tree, if it exists.
func (bb *BBTree) Remove(key Comparable) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	bb.root = bb.remove(key, bb.root)
	bb.root.parentColor = black
}

func (bb *BBTree) remove(key Comparable, currentNode *node) *node {
	if key.More(currentNode.key) {
		if !isBlue(currentNode.right) && !isBlue(currentNode.right.right) {
			currentNode = bb.moveBlueRight(currentNode)
		}
		currentNode.right = bb.remove(key, currentNode.right)
		return bb.fixUp(currentNode)
	}
	if isBlue(currentNode.right) {
		currentNode = bb.rotateLeft(currentNode)
	}

	if !isBlue(currentNode.left) && currentNode.left != nil && !isBlue(currentNode.left.right) {
		currentNode = bb.moveBlueLeft(currentNode)
	}

	// currentNode.key == key (By totality axiom.)
	if !currentNode.key.More(key) {
		if currentNode.left == nil {
			return nil
		}
		maxNode := bb.max(currentNode.left)
		currentNode.key = maxNode.Key
		currentNode.value = maxNode.Value
		currentNode.left = bb.removeMax(currentNode.left)
		return bb.fixUp(currentNode)
	}
	currentNode.left = bb.remove(key, currentNode.left)
	return bb.fixUp(currentNode)
}

// Find takes in a given Comparable key and returns a Node type and
// bool indicating whether the given key was found in the tree.
func (bb *BBTree) Find(key Comparable) (*Node, bool) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	return bb.find(key, bb.root)
}

func (bb *BBTree) find(key Comparable, n *node) (*Node, bool) {
	if n == nil {
		return nil, false
	}

	if key.More(n.key) {
		return bb.find(key, n.right)
	}
	if n.key.More(key) {
		return bb.find(key, n.left)
	}
	return &Node{n.key, n.value}, true
}

// Size returns the number of nodes in the tree.
func (bb *BBTree) Size() int {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	return size(bb.root)
}

// Rank returns the rank of a given key. Also known as the number
// of nodes in the tree less than the key.
func (bb *BBTree) Rank(key Comparable) int {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	return bb.rank(key, bb.root)
}

func (bb *BBTree) rank(key Comparable, currentNode *node) int {
	if currentNode == nil {
		return 0
	}

	if currentNode.key.More(key) {
		return bb.rank(key, currentNode.left)
	}

	if key.More(currentNode.key) {
		return size(currentNode.left) + 1 + bb.rank(key, currentNode.right)
	}

	// key == currentNode.key
	return size(currentNode.left)
}

// Select searches for the rank of a given size and returns the Node,
// if it exists.
func (bb *BBTree) Select(givenSize int) *Node {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	return bb.selector(givenSize, bb.root)
}

func (bb *BBTree) selector(givenSize int, currentNode *node) *Node {
	if currentNode == nil {
		return &Node{nil, nil}
	}

	if givenSize < size(currentNode.left) {
		return bb.selector(givenSize, currentNode.left)
	}

	if rank := size(currentNode.left); givenSize > rank {
		givenSize = givenSize - rank - 1
		return bb.selector(givenSize, currentNode.right)
	}

	return &Node{currentNode.key, currentNode.value}
}

// Count returns the number of nodes that is within an *inclusive* range
// of given low and high Comparable keys.
func (bb *BBTree) Count(lowKey, highKey Comparable) int {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	// Make sure range well formed.
	if lowKey.More(highKey) {
		return 0
	}

	count := bb.rank(highKey, bb.root) - bb.rank(lowKey, bb.root)
	if _, ok := bb.find(highKey, bb.root); ok {
		count++
	}
	return count
}

// Range returns a slice of Node types in sorted order within an
// *exclusive* range of given lower and upper bound Comparable keys.
func (bb *BBTree) Range(lowKey, highKey Comparable) []*Node {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	var keyValues []*Node
	bb.rangeSearch(&keyValues, lowKey, highKey, bb.root)
	return keyValues
}

// rangeSearch uses side-effects to build the keyValues slice instead
// of merging slices together.
func (bb *BBTree) rangeSearch(keyValues *[]*Node, lowKey Comparable, highKey Comparable, currentNode *node) {
	if currentNode == nil {
		return
	}

	if currentNode.key.More(lowKey) {
		bb.rangeSearch(keyValues, lowKey, highKey, currentNode.left)
	}

	if currentNode.key.More(lowKey) && highKey.More(currentNode.key) {
		*keyValues = append(*keyValues, &Node{currentNode.key, currentNode.value})
	}

	if highKey.More(currentNode.key) {
		bb.rangeSearch(keyValues, lowKey, highKey, currentNode.right)
	}
}

// ---- Sorting and Tree Representation methods

// InOrder returns a slice of all the Nodes in sorted order.
// Implements the InOrder recursive algorithm.
func (bb *BBTree) InOrder() []*Node {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	var sorted []*Node
	if bb.IsEmpty() {
		return sorted
	}

	bb.inOrder(&sorted, bb.root)
	return sorted
}

func (bb *BBTree) inOrder(sorted *[]*Node, currentNode *node) {
	if currentNode == nil {
		return
	}
	bb.inOrder(sorted, currentNode.left)
	*sorted = append(*sorted, &Node{currentNode.key, currentNode.value})
	bb.inOrder(sorted, currentNode.right)
}

// Print prints a representation of the tree in level order for debugging
// purposes.
func (bb *BBTree) Print() {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	bb.printify(bb.root, 0)
}

func (bb *BBTree) printify(currentNode *node, level int) {
	if currentNode == nil {
		return
	}
	format := strings.Repeat("    ", level)

	nodeString := "--[%v"
	if isBlue(currentNode) {
		nodeString = "==[%v"
	}
	toString := fmt.Sprintf(nodeString+"\n", currentNode)
	level++

	bb.printify(currentNode.left, level)
	fmt.Printf(format + toString)
	bb.printify(currentNode.right, level)
}

// New is the default constructor for a Black and Blue tree.
func New() *BBTree {
	return &BBTree{nil, sync.Mutex{}}
}
