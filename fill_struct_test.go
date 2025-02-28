package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFill_Simple(t *testing.T) {
	type testStruct struct {
		IntValue     int
		UintValue    uint
		FloatValue   float64
		ComplexValue complex64
		BoolValue    bool
		StringValue  string
		ArrayValue   [4]int
		SliceValue   []uint
		MapValue     map[string]float64
		// Can't do simple comparison of channel - excluded from this test
		//ChannelValue chan float64
	}

	expected := testStruct{
		IntValue:     -1,
		UintValue:    1,
		FloatValue:   1.234,
		ComplexValue: 1 + 2i,
		BoolValue:    true,
		StringValue:  "string",
		ArrayValue:   [4]int{-2, -3, -4, -5},
		SliceValue:   []uint{1, 1, 1, 1},
		MapValue:     map[string]float64{"string": 1.234},
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
	// IntValueField
	c.pushInt64(-1, NativeBytes)
	// ArrayValue elements
	c.pushInt64(-2, NativeBytes)
	c.pushInt64(-3, NativeBytes)
	c.pushInt64(-4, NativeBytes)
	c.pushInt64(-5, NativeBytes)

	return c
}

func TestFill_Channel(t *testing.T) {
	type testStruct struct {
		ChanValue chan float64
	}

	// Test value
	val := testStruct{}
	Fill(&val, nil)
	assert.Equal(t, 1, len(val.ChanValue))
	assert.Equal(t, 1.234, <-val.ChanValue)
}
