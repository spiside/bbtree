package bbtree

import (
	"math/rand"
	"testing"
)

var (
	expectedValues = map[Int]string{
		4: "hello",
		2: "world",
		3: "foo",
		1: "bar",
		5: "other",
	}
)

func isSorted(a []*Node) bool {
	for i := len(a) - 1; i > 0; i-- {
		if !a[i].Key.More(a[i-1].Key) {
			return false
		}
	}
	return true
}

func TestNewBBtree(t *testing.T) {
	bbtree := New()

	if bbtree == nil {
		t.Errorf("Constructor should create new BST")
	}
}

func TestPutWorks(t *testing.T) {
	bbtree := New()

	for i := 0; i < 40; i++ {
		bbtree.Put(Int(i), "")
	}
	sorted := bbtree.InOrder()
	if !isSorted(sorted) {
		t.Errorf("InOrder() does not return a sorted slice")
	}
}

func TestFindWorked(t *testing.T) {
	bbtree := New()
	rand.Seed(0)

	for _, k := range rand.Perm(40) {
		bbtree.Put(Int(k), "")
	}

	_, ok := bbtree.Find(Int(4))
	if !ok {
		t.Errorf("Find did not work as expected")
	}
}

func TestRemoveWorked(t *testing.T) {
	bbtree := New()

	rand.Seed(0)
	for _, k := range rand.Perm(1000) {
		bbtree.Put(Int(k), "")
	}

	bbtree.Remove(Int(4))

	_, ok := bbtree.Find(Int(4))
	if ok {
		t.Errorf("Remove did not work as expected")
	}

	sorted := bbtree.InOrder()
	if !isSorted(sorted) {
		t.Errorf("InOrder() does not return a sorted slice")
	}
}

func TestMinAndMaxWorkAsExpected(t *testing.T) {
	bbtree := New()

	rand.Seed(1)
	for _, k := range rand.Perm(100) {
		bbtree.Put(Int(k), "something")
	}

	minNode := bbtree.Min()
	maxNode := bbtree.Max()

	if minNode.Key != Int(0) {
		t.Errorf("Min not working as expected")
	}

	if maxNode.Key != Int(99) {
		t.Errorf("Max not working as expected")
	}
}

func TestRemoveMinWorkAsExpected(t *testing.T) {
	bbtree := New()

	rand.Seed(1)
	permutation := rand.Perm(100)
	for _, k := range permutation {
		bbtree.Put(Int(k), "")
	}

	bbtree.RemoveMin()
	if _, ok := bbtree.Find(Int(0)); ok {
		t.Errorf("RemoveMin not working as expected")
	}

	for i := 1; i <= 19; i++ {
		bbtree.RemoveMin()
	}

	if count := bbtree.Count(Int(-1), Int(19)); count > 0 {
		t.Errorf("Repeated calls to RemoveMin not working as expected")
	}
}

func TestRemoveMaxWorksAsExpected(t *testing.T) {
	bbtree := New()

	rand.Seed(3)
	permutation := rand.Perm(100)
	for _, k := range permutation {
		bbtree.Put(Int(k), "")
	}

	bbtree.RemoveMax()
	if _, ok := bbtree.Find(Int(99)); ok {
		t.Errorf("RemoveMin not working as expected")
	}

	for i := 1; i <= 19; i++ {
		bbtree.RemoveMax()
	}

	if count := bbtree.Count(Int(80), Int(99)); count > 0 {
		t.Errorf("Repeated calls to RemoveMin not working as expected")
	}
}

func TestSizeWorks(t *testing.T) {
	bbtree := New()

	rand.Seed(2)
	permutation := rand.Perm(100)
	for _, k := range permutation {
		bbtree.Put(Int(k), "something")
	}

	if size := bbtree.Size(); size != 100 {
		t.Errorf("Size() returns wrong value of %v when adding", size)
	}

	for _, k := range permutation[:10] {
		bbtree.Remove(Int(k))
	}

	if size := bbtree.Size(); size != 90 {
		t.Errorf("Size() returns wrong value of %v when removing", size)
	}
}

func TestRangeWorks(t *testing.T) {
	bbtree := New()

	rand.Seed(5)
	permutation := rand.Perm(100)
	for _, k := range permutation {
		bbtree.Put(Int(k), "")
	}

	ranged := bbtree.Range(Int(-1), Int(100))
	if len(ranged) != 100 {
		t.Errorf("Range() does not return a slice of expected Nodes")
	}

	if !isSorted(ranged) {
		t.Errorf("Range() returns slice in unsorted order")
	}

	ranged = bbtree.Range(Int(10), Int(21))
	if len(ranged) != 10 {
		t.Errorf("Range() does not return a slice of expected Nodes")
	}

	if !isSorted(ranged) {
		t.Errorf("Range() returns slice in unsorted order")
	}

}

func TestCountWorks(t *testing.T) {
	bbtree := New()

	rand.Seed(2)
	permutation := rand.Perm(100)
	for _, k := range permutation {
		bbtree.Put(Int(k), "something")
	}

	if count := bbtree.Count(Int(0), Int(99)); count != 100 {
		t.Errorf("Count not working as expected")
	}

	if count := bbtree.Count(Int(11), Int(20)); count != 10 {
		t.Errorf("Count not working as expected")
	}

}
