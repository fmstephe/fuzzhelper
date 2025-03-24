package fuzzhelper

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestFill_SimpleTypes(t *testing.T) {
	type testStruct struct {
		IntValue   int
		Int64Value int64
		Int32Value int32
		Int16Value int16
		Int8Value  int8

		UintValue   uint
		Uint64Value uint64
		Uint32Value uint32
		Uint16Value uint16
		Uint8Value  uint8

		Float64Value float64
		Float32Value float32

		Bool1Value bool
		Bool2Value bool

		String1Value string
		String2Value string
		String3Value string
		String4Value string

		ArrayValue [4]int
		SliceValue []uint

		MapValue map[string]float64
	}

	expected := testStruct{
		IntValue:   -1,
		Int64Value: -64,
		Int32Value: -32,
		Int16Value: -16,
		Int8Value:  -8,

		UintValue:   1,
		Uint64Value: 64,
		Uint32Value: 32,
		Uint16Value: 16,
		Uint8Value:  8,

		Float64Value: 3.1415,
		Float32Value: 4.1415,

		Bool1Value: true,
		Bool2Value: false,

		String1Value: "a",
		String2Value: "ab",
		String3Value: "abc",
		String4Value: "abcd",

		ArrayValue: [4]int{-2, -3, -4, -5},
		SliceValue: []uint{2, 3, 4, 5},

		MapValue: map[string]float64{"map key string": 5.1415},
	}

	// Test value
	val := testStruct{}
	Fill(&val, buildSimpleTestByteConsumer())

	assert.Equal(t, expected, val)

	// Test pointer
	valp := &testStruct{}
	Fill(valp, buildSimpleTestByteConsumer())
	assert.Equal(t, expected, *valp)

	// Test pointer to pointer
	var valpp *testStruct
	Fill(&valpp, buildSimpleTestByteConsumer())
	assert.NotNil(t, valpp)
	assert.Equal(t, expected, *valpp)
}

// Builds a ByteConsumer which should build a testStruct which matches the expected values
func buildSimpleTestByteConsumer() *ByteConsumer {
	// Set all the fill values here
	c := NewByteConsumer([]byte{})

	// IntValue fields
	c.pushInt64(-1, BytesForNative)
	c.pushInt64(-64, BytesFor64)
	c.pushInt64(-32, BytesFor32)
	c.pushInt64(-16, BytesFor16)
	c.pushInt64(-8, BytesFor8)

	// UintValue field
	c.pushUint64(1, BytesForNative)
	c.pushUint64(64, BytesFor64)
	c.pushUint64(32, BytesFor32)
	c.pushUint64(16, BytesFor16)
	c.pushUint64(8, BytesFor8)

	// Float64Value field
	c.pushFloat64(3.1415, BytesFor64)
	// Float32Value field
	c.pushFloat64(4.1415, BytesFor32)

	// Bool1Value field
	c.pushBool(true)
	// Bool2Value field
	c.pushBool(false)

	// StringField field
	c.pushString("a")
	c.pushString("ab")
	c.pushString("abc")
	c.pushString("abcd")

	// ArrayValue has fixed size, requires no data
	// ArrayValue Elements
	c.pushInt64(-2, BytesForNative)
	c.pushInt64(-3, BytesForNative)
	c.pushInt64(-4, BytesForNative)
	c.pushInt64(-5, BytesForNative)

	// SliceValue Size
	c.pushUint64(4, BytesForNative)
	// SliceValue Elements
	c.pushUint64(2, BytesForNative)
	c.pushUint64(3, BytesForNative)
	c.pushUint64(4, BytesForNative)
	c.pushUint64(5, BytesForNative)

	// Map Size
	c.pushInt64(1, BytesForNative)
	// MapValue map key
	c.pushString("map key string")
	// MapValue map entry
	c.pushFloat64(5.1415, BytesFor64)

	return c
}

func TestFill_Map(t *testing.T) {
	type valueStruct struct {
		IntField int
	}
	type testStruct struct {
		MapValue map[int]valueStruct
	}

	// Set all the fill values here
	c := NewByteConsumer([]byte{})
	// map is size 1
	c.pushUint64(1, BytesForNative)
	// IntValue field
	c.pushUint64(1, BytesForNative)
	c.pushUint64(2, BytesForNative)

	// Test value
	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, 1, len(val.MapValue))
	assert.Equal(t, 2, val.MapValue[1].IntField)
}

// Test a series of nested structs.
//
// This test is very very complex, hard to write - harder to read. It's not
// clear if it's a good test. But we do need to test this case somehow. We'll
// have a think about this. Maybe the test will prove resilient and not need to
// be changed in the future, in that case it will probably be left alone. But
// if it turns out to be fragile and needs constant changes we will revisit it.
func TestFill_Complex(t *testing.T) {
	type innerInnerStruct struct {
		UintValue   uint
		StringValue string
	}

	innerInnerF := func() innerInnerStruct {
		return innerInnerStruct{
			UintValue:   1,
			StringValue: "innerinner",
		}
	}

	innerInnerBytesF := func(c *ByteConsumer) {
		c.pushInt64(1, BytesForNative)
		c.pushString("innerinner")
	}

	type innerStruct struct {
		IntValue    int
		InnerInnerP innerInnerStruct
		StringValue string
	}

	innerF := func() innerStruct {
		innerInnerP := innerInnerF()
		return innerStruct{
			IntValue:    -2,
			InnerInnerP: innerInnerP,
			StringValue: "inner",
		}
	}

	innerBytesF := func(c *ByteConsumer) {
		c.pushInt64(-2, BytesForNative)
		innerInnerBytesF(c)
		c.pushString("inner")
	}

	type testStruct struct {
		InnerV   innerStruct
		MapField map[string]innerStruct
	}

	c := NewByteConsumer([]byte{})
	// First layer of InnerV field
	innerBytesF(c)
	// Map size
	c.pushInt64(1, BytesForNative)
	c.pushString("key")
	innerBytesF(c)

	inner := innerF()
	expected := testStruct{
		InnerV: inner,
		MapField: map[string]innerStruct{
			"key": inner,
		},
	}

	// Test value
	val := testStruct{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

func TestLinkedList_One(t *testing.T) {
	type node struct {
		Value int
		Next  *node
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)

	expected := node{
		Value: 1,
		Next:  nil, // <- ran out of data here
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

func TestLinkedList_Two(t *testing.T) {
	type node struct {
		Value int
		Next  *node
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)

	expected := node{
		Value: 1,
		Next: &node{
			Value: 2,
			Next:  nil, // <- ran out of data here
		},
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

func TestLinkedList_Three(t *testing.T) {
	type node struct {
		Value int
		Next  *node
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)

	expected := node{
		Value: 1,
		Next: &node{
			Value: 2,
			Next: &node{
				Value: 3,
				Next:  nil, // <- ran out of data here
			},
		},
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

// NB: This test clarifies that trying to build a binary tree produces a
// roughly balanced tree.
func TestBalancedBinaryTree(t *testing.T) {
	type node struct {
		Value      int
		LeftChild  *node
		RightChild *node
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)
	c.pushInt64(5, BytesForNative)
	c.pushInt64(6, BytesForNative)

	expected := node{
		Value: 1,
		LeftChild: &node{
			Value: 2,
			LeftChild: &node{
				Value:      4,
				LeftChild:  &node{},
				RightChild: &node{},
			},
			RightChild: &node{
				Value:      5,
				LeftChild:  &node{},
				RightChild: &node{},
			},
		},
		RightChild: &node{
			Value: 3,
			LeftChild: &node{
				Value:      6,
				LeftChild:  nil, // <- ran out of data here
				RightChild: nil,
			},
			RightChild: &node{
				Value:      0,
				LeftChild:  nil,
				RightChild: nil,
			},
		},
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

func TestFill_UnsupportedTypes(t *testing.T) {
	type testStruct struct {
		ChanField          chan int
		InterfaceField     any
		ComplexField       complex128
		FuncField          func()
		UintptrField       uintptr
		UnsafePointerField unsafe.Pointer
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)
	c.pushInt64(5, BytesForNative)
	c.pushInt64(6, BytesForNative)
	c.pushInt64(7, BytesForNative)
	c.pushInt64(8, BytesForNative)
	c.pushInt64(9, BytesForNative)
	c.pushInt64(10, BytesForNative)
	c.pushInt64(11, BytesForNative)
	c.pushInt64(12, BytesForNative)

	val := &testStruct{}
	Fill(val, c)

	// Assert that none of those fields are set
	assert.Equal(t, &testStruct{}, val)
}

// Demonstrate that when the root of the value passed into Fill() is a pointer to a slice
// Then the slice is appended to until all of the data is used up
func TestFill_RootSlice(t *testing.T) {
	type testStruct struct {
		IntField int
	}

	{
		// Create root slice with enough data for one element
		c := NewByteConsumer([]byte{})
		c.pushInt64(1, BytesForNative)

		val := &[]testStruct{}
		Fill(val, c)

		assert.Equal(t, &[]testStruct{
			testStruct{1},
		}, val)
	}

	{
		// Create root slice with enough data for two elements
		c := NewByteConsumer([]byte{})
		c.pushInt64(1, BytesForNative)
		c.pushInt64(2, BytesForNative)

		val := &[]testStruct{}
		Fill(val, c)

		assert.Equal(t, &[]testStruct{
			testStruct{1},
			testStruct{2},
		}, val)
	}

	{
		// Create root slice with enough data for three elements
		c := NewByteConsumer([]byte{})
		c.pushInt64(1, BytesForNative)
		c.pushInt64(2, BytesForNative)
		c.pushInt64(3, BytesForNative)

		val := &[]testStruct{}
		Fill(val, c)

		assert.Equal(t, &[]testStruct{
			testStruct{1},
			testStruct{2},
			testStruct{3},
		}, val)
	}

	{
		// Create root slice with enough data for twelve elements
		c := NewByteConsumer([]byte{})
		c.pushInt64(1, BytesForNative)
		c.pushInt64(2, BytesForNative)
		c.pushInt64(3, BytesForNative)
		c.pushInt64(4, BytesForNative)
		c.pushInt64(5, BytesForNative)
		c.pushInt64(6, BytesForNative)
		c.pushInt64(7, BytesForNative)
		c.pushInt64(8, BytesForNative)
		c.pushInt64(9, BytesForNative)
		c.pushInt64(10, BytesForNative)
		c.pushInt64(11, BytesForNative)
		c.pushInt64(12, BytesForNative)

		val := &[]testStruct{}
		Fill(val, c)

		// Assert that none of those fields are set
		assert.Equal(t, &[]testStruct{
			testStruct{1},
			testStruct{2},
			testStruct{3},
			testStruct{4},
			testStruct{5},
			testStruct{6},
			testStruct{7},
			testStruct{8},
			testStruct{9},
			testStruct{10},
			testStruct{11},
			testStruct{12},
		}, val)
	}
}
