package sortedtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushOne(t *testing.T) {
	tree := New()

	tree.Add("a")

	assert.Equal(t, "a", tree.Least())
	assert.Equal(t, "a", tree.Greatest())
}

func TestPushTwo(t *testing.T) {
	tree := New()

	tree.Add("a")
	tree.Add("b")

	assert.Equal(t, "a", tree.Least())
	assert.Equal(t, "b", tree.Greatest())
}

func TestPushThree(t *testing.T) {
	tree := New()

	tree.Add("a")
	tree.Add("b")
	tree.Add("c")

	assert.Equal(t, "a", tree.Least())
	assert.Equal(t, "c", tree.Greatest())
}

func TestPushThreeReverse(t *testing.T) {
	tree := New()

	tree.Add("c")
	tree.Add("b")
	tree.Add("a")

	assert.Equal(t, "a", tree.Least())
	assert.Equal(t, "c", tree.Greatest())
}

func TestPushManyMixed(t *testing.T) {
	tree := New()

	tree.Add("c")
	tree.Add("b")
	tree.Add("a")
	tree.Add("a")
	tree.Add("b")
	tree.Add("b")
	tree.Add("z")

	assert.Equal(t, "a", tree.Least())
	assert.Equal(t, "z", tree.Greatest())
}
