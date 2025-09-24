package tests

import (
	"fmt"
	"github.com/PRITISH-TOMAR/byted/internal/btree"
	"testing"
)

func TestBPlusTree(t *testing.T) {
	// create tree with order 4 (small for easy splitting in test)
	tree := btree.New(4)

	// insert some keys
	for i := 1; i <= 21; i++ {
		key := fmt.Sprintf("k%02d", i)
		val := fmt.Sprintf("val-%02d", i)
		tree.Insert(key, val)
	}
	t.Log("Tree structure:\n", tree.Print())

	// exact lookup
	for i := 1; i <= 21; i++ {
		key := fmt.Sprintf("k%02d", i)
		v, ok := tree.Get(key)
		fmt.Printf("Get %s -> %v, %v\n", key, v, ok)
	}

	// range query k02..k06
	// fmt.Println("Range k02..k06:")
	// res := t.Range("k02", "k06", 0)
	// for _, p := range res {
	// 	fmt.Printf("  %s -> %v\n", p.Key, p.Value)
	// }

	// delete a key
	// ok = t.Delete("k05")
	// fmt.Println("Deleted k05?", ok)
	// fmt.Println("After delete:\n", t.DebugString())

	// range all
	// all := t.Range("k00", "k99", 0)
	// fmt.Println("All keys:")
	// for _, p := range all {
	// 	fmt.Printf("  %s -> %v\n", p.Key, p.Value)
	// }
}
