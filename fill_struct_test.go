package fuzzhelper

import (
	"testing"

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

		ComplexValue complex64

		Bool1Value bool
		Bool2Value bool

		String1Value string
		String2Value string
		String3Value string
		String4Value string

		ArrayValue [4]int
		SliceValue []uint

		MapValue map[string]float64
		// Can't do simple comparison of channel
		//ChannelValue chan float64
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

		ComplexValue: 1 + 2i,

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

	// ArrayValue elements
	c.pushInt64(-2, BytesForNative)
	c.pushInt64(-3, BytesForNative)
	c.pushInt64(-4, BytesForNative)
	c.pushInt64(-5, BytesForNative)

	// SliceValue elements
	c.pushUint64(4, BytesForNative)
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

func TestFill_Channel(t *testing.T) {
	type testStruct struct {
		ChanValue chan float64
	}

	// Set all the fill values here
	c := NewByteConsumer([]byte{})
	// Channel is size 1
	c.pushUint64(1, BytesForNative)
	// IntValue field
	c.pushFloat64(3.1415, BytesFor64)

	// Test value
	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, 1, len(val.ChanValue))
	assert.Equal(t, 3.1415, <-val.ChanValue)
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
		IntValue     int
		UintValue    uint
		Float64Value float64
		StringValue  string
	}

	innerInnerF := func() innerInnerStruct {
		return innerInnerStruct{
			IntValue:     -1,
			UintValue:    1,
			Float64Value: 1.234,
			StringValue:  "string",
		}
	}

	innerInnerBytesF := func(c *ByteConsumer) {
		c.pushInt64(-1, BytesForNative)
		c.pushUint64(1, BytesForNative)
		c.pushFloat64(1.234, BytesFor64)
		c.pushString("string")
	}

	type innerStruct struct {
		IntValue     int
		InnerInnerP  *innerInnerStruct
		UintValue    uint
		InnerInnerV  innerInnerStruct
		Float64Value float64
	}

	innerF := func() innerStruct {
		innerInnerP := innerInnerF()
		return innerStruct{
			IntValue:     -2,
			InnerInnerP:  &innerInnerP,
			UintValue:    2,
			InnerInnerV:  innerInnerF(),
			Float64Value: 2.234,
		}
	}

	innerBytesF := func(c *ByteConsumer) {
		c.pushInt64(-2, BytesForNative)
		innerInnerBytesF(c)
		c.pushUint64(2, BytesForNative)
		innerInnerBytesF(c)
		c.pushFloat64(2.234, BytesFor64)
	}

	type testStruct struct {
		InnerPP  **innerStruct
		InnerP   *innerStruct
		InnerV   innerStruct
		MapField map[string]innerStruct
	}

	c := NewByteConsumer([]byte{})
	innerBytesF(c)
	innerBytesF(c)
	innerBytesF(c)
	c.pushInt64(1, BytesForNative)
	c.pushString("key")
	innerBytesF(c)

	inner := innerF()
	innerP := &inner
	innerPP := &innerP
	expected := testStruct{
		InnerPP: innerPP,
		InnerP:  innerP,
		InnerV:  inner,
		MapField: map[string]innerStruct{
			"key": innerF(),
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
		Next:  nil,
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
			Next:  nil,
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
				Next:  nil,
			},
		},
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}

// NB: This test clarifies that trying to build a binary tree produces a
// lefthanded fully unbalanced tree.
func TestUnbalancedBinaryTree(t *testing.T) {
	type node struct {
		Value      int
		LeftChild  *node
		RightChild *node
	}

	c := NewByteConsumer([]byte{})
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)

	expected := node{
		Value: 1,
		LeftChild: &node{
			Value: 2,
			LeftChild: &node{
				Value:      3,
				LeftChild:  nil,
				RightChild: nil,
			},
			RightChild: nil,
		},
		RightChild: nil,
	}

	val := node{}
	Fill(&val, c)

	assert.Equal(t, expected, val)
}
