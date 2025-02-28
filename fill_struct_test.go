package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFill_int(t *testing.T) {
	type testStruct struct {
		IntValue    int
		UintValue   uint
		FloatValue  float64
		BoolValue   bool
		StringValue string
		ArrayValue  [4]int
		SliceValue  []uint
	}

	expected := testStruct{
		IntValue:    -1,
		UintValue:   1,
		FloatValue:  1.234,
		BoolValue:   true,
		StringValue: "string",
		ArrayValue:  [4]int{-1, -1, -1, -1},
		SliceValue:  []uint{1, 1, 1, 1},
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
