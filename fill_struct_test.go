package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFill_int(t *testing.T) {
	type intstruct struct {
		Value int
	}

	// Test value
	val := intstruct{}
	Fill(&val)
	assert.Equal(t, -1, val.Value)

	// Test pointer
	valp := &intstruct{}
	Fill(valp)
	assert.Equal(t, -1, valp.Value)

	// Test pointer to pointer
	var valpp *intstruct
	Fill(&valpp)
	assert.NotNil(t, valpp)
	assert.Equal(t, -1, valpp.Value)
}

func TestFill_uint(t *testing.T) {
	type uintstruct struct {
		Value uint
	}

	// Test value
	val := uintstruct{}
	Fill(&val)
	assert.Equal(t, uint(1), val.Value)

	// Test pointer
	valp := &uintstruct{}
	Fill(valp)
	assert.Equal(t, uint(1), valp.Value)

	// Test pointer to pointer
	var valpp *uintstruct
	Fill(&valpp)
	assert.NotNil(t, valpp)
	assert.Equal(t, uint(1), valpp.Value)
}

func TestFill_float64(t *testing.T) {
	type float64struct struct {
		Value float64
	}

	// Test value
	val := float64struct{}
	Fill(&val)
	assert.Equal(t, float64(1.234), val.Value)

	// Test pointer
	valp := &float64struct{}
	Fill(valp)
	assert.Equal(t, float64(1.234), valp.Value)

	// Test pointer ot pointer
	var valpp *float64struct
	Fill(&valpp)
	assert.NotNil(t, &valpp)
	assert.Equal(t, float64(1.234), valpp.Value)
}

func TestFill_bool(t *testing.T) {
	type boolstruct struct {
		Value bool
	}

	// Test value
	val := boolstruct{}
	Fill(&val)
	assert.Equal(t, true, val.Value)

	// Test pointer
	valp := &boolstruct{}
	Fill(valp)
	assert.Equal(t, true, valp.Value)

	// Test pointer ot pointer
	var valpp *boolstruct
	Fill(&valpp)
	assert.NotNil(t, &valpp)
	assert.Equal(t, true, valpp.Value)
}

func TestFill_string(t *testing.T) {
	type stringstruct struct {
		Value string
	}

	// Test value
	val := stringstruct{}
	Fill(&val)
	assert.Equal(t, "string", val.Value)

	// Test pointer
	valp := &stringstruct{}
	Fill(valp)
	assert.Equal(t, "string", valp.Value)

	// Test pointer ot pointer
	var valpp *stringstruct
	Fill(&valpp)
	assert.NotNil(t, &valpp)
	assert.Equal(t, "string", valpp.Value)
}
