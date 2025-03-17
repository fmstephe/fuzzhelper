package sortedtree

import "fmt"

// A tree for sorting strings. It's not self balancing, and is as simple as I
// could manage in its implementation.
//
// It has no practical uses and only exists to demonstrate a bug which would be
// hard to catch with unit tests but easy to catch with fuzzing.
//
// The bug itself can be found in the comparison function which only compares
// the first 127 characters of the strings being sorted. This will produce bad
// behaviour when sorting longer strings.
type SortedTree struct {
	size int
	root *node
}

func New() *SortedTree {
	return &SortedTree{}
}

func (t *SortedTree) Add(s string) {
	t.root = t.root.add(s)
	t.size++
}

func (t *SortedTree) Least() string {
	return t.root.least()
}

func (t *SortedTree) Greatest() string {
	return t.root.greatest()
}

func (t *SortedTree) String() string {
	return fmt.Sprintf("size: %d, least: %s, greatest: %s", t.size, t.Least(), t.Greatest())
}

type node struct {
	smaller *node
	larger  *node
	s       string
}

func (n *node) least() string {
	if n.smaller == nil {
		return n.s
	}
	return n.smaller.least()
}

func (n *node) greatest() string {
	if n.larger == nil {
		return n.s
	}
	return n.larger.greatest()
}

func (n *node) add(s string) *node {
	if n == nil {
		// If this method is called on a nil node
		// Create a new leaf
		return &node{
			s: s,
		}
	}

	if less(s, n.s) {
		n.smaller = n.smaller.add(s)
		return n
	}

	n.larger = n.larger.add(s)
	return n
}

// This function contains an obvious bug, it only compares the beginning of
// long strings. This limit will make the comparison misbehave with some
// strings, but for most conventional unit tests the function will appear to be
// working.
//
// Fuzzing should find this bug
func less(s1, s2 string) bool {
	s1 = limit(s1)
	s2 = limit(s2)

	return s1 < s2
}

// shorten long strings - this is a bug, it means that long strongs are not compared properly
func limit(s string) string {
	const limit = 64
	if len(s) > limit {
		return s[:limit]
	}
	return s
}
