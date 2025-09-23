package btree

import (
	"fmt"
	"sort"
)

// BPlusTree represents a B+ tree data structure.

// Core Idea :
// 1. keys are strings.lexicographical order
// 2. values are stored in leaf nodes as pointers to any valueMeta.
// 3. leaf nodes linkdby `next` pointer for range queries.
// 4. internal nodes store only keys and child pointers.
// 5. all leaf nodes are at the same level.
// 6. order :  max number of keys for internal nodes.

// For this phase : deletion  : NOT BALANCING B+ TREE

const defaultOrder = 32

// Key-Value Pair Structure
type KVPair struct {
	Key   string
	Value any // can be pointer to valueMeta or actual value.
}

// Node structure
type Node struct {
	isLeaf   bool
	keys     []string
	children []*Node // for internal nodes: child pointers
	values   []any   // for leaf nodes: values corresponding to keys
	next     *Node   // pointer to next leaf node (for leaf nodes only)
	parent   *Node   // pointer to parent node (nil for root)
}

// B+ Tree main structure
type BPlusTree struct {
	root  *Node
	order int // max number of keys in internal nodes
}

// New creates a new B+ tree with the specified order.
func New(order int) *BPlusTree {
	if order <= 0 {
		order = defaultOrder
	}
	// Initialize root as a leaf node
	res := &Node{
		isLeaf:   true,
		keys:     make([]string, 0),
		values:   make([]any, 0),
		children: nil,
		next:     nil,
		parent:   nil,
	}

	return &BPlusTree{root: res, order: order}
}

// find for a certain key in the tree
func (t *BPlusTree) findLeaf(key string) *Node {
	curr := t.root

	for !curr.isLeaf {
		// find the child pointer to follow
		i := sort.SearchStrings(curr.keys, key)

		curr = curr.children[i]
	}

	return curr
}


func (t * BPlusTree) Get(key string) (any, bool) {
	leaf 
}