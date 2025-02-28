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
		ArrayValue:   [4]int{-1, -1, -1, -1},
		SliceValue:   []uint{1, 1, 1, 1},
		MapValue:     map[string]float64{"string": 1.234},
	}

	// Test value
	val := testStruct{}
	Fill(&val)
	assert.Equal(t, expected, val)

	// Test pointer
	valp := &testStruct{}
	Fill(valp)
	assert.Equal(t, expected, *valp)

	// Test pointer to pointer
	var valpp *testStruct
	Fill(&valpp)
	assert.NotNil(t, valpp)
	assert.Equal(t, expected, *valpp)
}

func TestFill_Channel(t *testing.T) {
	type testStruct struct {
		ChanValue chan float64
	}

	// Test value
	val := testStruct{}
	Fill(&val)
	assert.Equal(t, 1, len(val.ChanValue))
	assert.Equal(t, 1.234, <-val.ChanValue)
}
