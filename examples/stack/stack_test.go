package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_PopEmpty(t *testing.T) {
	s := New()

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())
}

func TestStack_PushPop(t *testing.T) {
	s := New()

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())

	s.Push("foo")
	val := s.Pop()
	assert.Equal(t, "foo", val)

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())
}

func TestStack_PushPushPopPop(t *testing.T) {
	s := New()

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())

	s.Push("foo")
	s.Push("bar")
	assert.Equal(t, "bar", s.Pop())
	assert.Equal(t, "foo", s.Pop())

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())
}

func TestStack_PushPushPopPushPopPop(t *testing.T) {
	s := New()

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())

	s.Push("foo")
	s.Push("bar")
	assert.Equal(t, "bar", s.Pop())
	s.Push("bang")
	assert.Equal(t, "bang", s.Pop())
	assert.Equal(t, "foo", s.Pop())

	// Popping when the stack is empty produces nil
	assert.Nil(t, s.Pop())
}
