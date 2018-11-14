package bbtree_test

import (
	"fmt"

	"github.com/spiside/bbtree"
)

func ExampleNew() {
	bb := bbtree.New()
	bb.Put(bbtree.Int(1), "hello")
	bb.Put(bbtree.Int(2), "world")

	_, ok := bb.Find(bbtree.Int(2))
	if ok {
		fmt.Println("Found key!")
	}
	// Output: Found key!
}
