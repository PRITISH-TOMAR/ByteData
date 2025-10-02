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

// find for a certain key's leaf node in the tree
func (t *BPlusTree) findLeaf(key string) *Node {
	curr := t.root

	for !curr.isLeaf {
		// Use sort.Search to find the child index
		i := sort.Search(len(curr.keys), func(i int) bool {
			return key < curr.keys[i]
		})
		curr = curr.children[i]
	}

	return curr
}

func (t *BPlusTree) Get(key string) (any, bool) {
	leaf := t.findLeaf(key)
	i := sort.SearchStrings(leaf.keys, key) // search for the key in the leaf node
	if i < len(leaf.keys) && leaf.keys[i] == key {
		return leaf.values[i], true // if found
	}
	return nil, false // if not found -> returns  nil, false
}

// InsertIntoLeaf inserts a key-value pair into a leaf node.
func (t *BPlusTree) InsertIntoLeaf(leaf *Node, key string, value any) {
	i := sort.SearchStrings(leaf.keys, key) // search for key

	// if key exists, replace
	if i < len(leaf.keys) && leaf.keys[i] == key {
		leaf.values[i] = value
		return
	}

	// insert at position i..
	leaf.keys = append(leaf.keys, "")      // make space for new key
	leaf.values = append(leaf.values, nil) // make space for new value

	copy(leaf.keys[i+1:], leaf.keys[i:]) // shift keys to the right
	leaf.keys[i] = key

	copy(leaf.values[i+1:], leaf.values[i:]) // shift values to the right
	leaf.values[i] = value
}

// slpitLeaf splits a leaf node and returns the new leaf node and the promoted key.
func (t *BPlusTree) splitLeaf(leaf *Node) (*Node, string) {
	mid := len(leaf.keys) / 2 // finding mid for balancing the split

	// Create new leaf node
	newRightLeaf := &Node{
		isLeaf: true,
		keys:   append([]string{}, leaf.keys[mid:]...), // copy second half of keys
		values: append([]any{}, leaf.values[mid:]...),  // copy second half of values
		next:   leaf.next,                              // new leaf points to the next of current leaf
		parent: leaf.parent,
	}

	// Update current leaf
	leaf.keys = leaf.keys[:mid] // from 0 to mid-1
	leaf.values = leaf.values[:mid]
	leaf.next = newRightLeaf // current leaf points to new leaf

	// Insert new key into parent
	promoted_key := newRightLeaf.keys[0] // first key of new leaf to be promoted
	return newRightLeaf, promoted_key    // return new leaf and its first key to be inserted into parent
}

// splitInternal splits an internal node and returns the new node and the promoted key.
func (t *BPlusTree) splitInternal(node *Node) (*Node, string) {
	mid := len(node.keys) / 2
	promoted_key := node.keys[mid] // first key of new internal to be promoted
	// Create new internal node
	newRightInternal := &Node{
		isLeaf:   false,
		keys:     append([]string{}, node.keys[mid+1:]...),    // copy second half of keys
		children: append([]*Node{}, node.children[mid+1:]...), // copy second half of children
		parent:   node.parent,
	}

	// Update parent pointers of moved children
	for _, child := range newRightInternal.children {
		if child != nil {
			child.parent = newRightInternal
		}
	}

	// Update current node -> to left node
	node.keys = node.keys[:mid] // from 0 to mid-1
	node.children = node.children[:mid+1]

	return newRightInternal, promoted_key
}

// InsertIntoParent inserts key and child into parent node.
func (t *BPlusTree) InsertIntoParent(left *Node, key string, right *Node) {
	parent := left.parent
	// if left is root
	if parent == nil {
		// create new root
		newRoot := &Node{
			isLeaf:   false,
			keys:     []string{key},
			children: []*Node{left, right},
			values:   nil,
			next:     nil,
			parent:   nil,
		}
		left.parent = newRoot
		right.parent = newRoot
		t.root = newRoot
		return
	}

	// insert key and right child into parent
	i := sort.SearchStrings(parent.keys, key)

	parent.keys = append(parent.keys, "")    // make space for new key
	copy(parent.keys[i+1:], parent.keys[i:]) // shift keys to the right
	parent.keys[i] = key

	parent.children = append(parent.children, nil)     // make space for new child
	copy(parent.children[i+2:], parent.children[i+1:]) // shift children to the right
	parent.children[i+1] = right                       // as left child is already at its correct position in children

	right.parent = parent

	// if parent overflows, split it
	if len(parent.keys) > t.order {
		// split intervals.
		rightSibling, promoted_key := t.splitInternal(parent)
		// recursively insert into parent
		t.InsertIntoParent(parent, promoted_key, rightSibling)
	}
}

// Insert inserts a key-value pair into the B+ tree.
func (t *BPlusTree) Insert(key string, value any) {
	if t.root == nil {
		t.root = &Node{
			isLeaf: true,
			keys:   []string{key},
			values: []any{value},
		}
		return
	}

	leaf := t.findLeaf(key)

	// insert or replace in leaf
	t.InsertIntoLeaf(leaf, key, value)

	// if leaf overflows, split it
	if len(leaf.keys) > t.order {
		newRightLeaf, promoted_key := t.splitLeaf(leaf)
		// recursively insert into parent

		t.InsertIntoParent(leaf, promoted_key, newRightLeaf)
	}
}

// Debugging or printing the tree structure
func (t *BPlusTree) Print() string {
	if t.root == nil {
		return "<EMPTY>"
	}
	// BFS traversal per level

	currentLevel := []*Node{t.root}
	result := ""
	level := 0

	for len(currentLevel) > 0 {
		result += fmt.Sprintf("Level %d: ", level)
		nextLevel := []*Node{}

		for _, n := range currentLevel {
			if n.isLeaf {
				result += fmt.Sprintf("Leaf keys: %v) | ", n.keys)
			} else {
				result += fmt.Sprintf("Internal keys: %v) | ", n.keys)
				nextLevel = append(nextLevel, n.children...)
			}
		}

		result += "\n"
		currentLevel = nextLevel
		level++
	}
	return result
}

// Delete removes a key from the B+ tree -> Phase 1 : No rebalancing .
func (t *BPlusTree) Delete(key string) bool {
	leaf := t.findLeaf(key)

	if leaf == nil {
		return false // tree is empty
	}

	i := sort.SearchStrings(leaf.keys, key) // search for the key in the leaf node
	if i < len(leaf.keys) && leaf.keys[i] == key {
		// Key found, remove it
		leaf.keys = append(leaf.keys[:i], leaf.keys[i+1:]...)
		leaf.values = append(leaf.values[:i], leaf.values[i+1:]...)
		// Note: Balancing after deletion ->  Planned for  upcoming version/phase 1+.
		return true
	}
	return false // key not found
}

// RangeQueries for range [start, end] -> actual purpose of B+ Tree
func (t *BPlusTree) RangeQuery(start, end string) []KVPair {
	results := []KVPair{}
	leaf := t.findLeaf(start)

	// 1. find the starting leaf node
	for leaf != nil {
		for i, key := range leaf.keys {
			if key >= start && key <= end {
				results = append(results, KVPair{Key: key, Value: leaf.values[i]})
			}
			if key > end {
				return results
			}
		}
		leaf = leaf.next // move to the next leaf node
	}

	return results
}
