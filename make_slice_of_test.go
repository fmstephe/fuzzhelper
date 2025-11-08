package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSliceOf(t *testing.T) {
	allowableTypes := []any{
		// value receiver method, value type
		&interfaceDemoA{},
		// pointer receiver method, pointer type
		&interfaceDemoB{},
		// value receiver method, pointer type
		&interfaceDemoC{},
		// pointer receiver method, value type
		&interfaceDemoD{},
		// value receiver method, value type
		&interfaceDemoE{},
		// pointer receiver method, pointer type
		&interfaceDemoF{},
	}
	c := newByteConsumer([]byte{})
	// Values for first two interface choices in slice
	c.pushUint64(0, bytesForNative)
	c.pushUint64(1, bytesForNative)

	// Values for interfaceExampleA
	c.pushInt64(-100, bytesForNative)
	c.pushFloat64(123.123, bytesFor64)

	// Continue pushing choices for the rest of the interfaces
	c.pushUint64(2, bytesForNative)
	c.pushUint64(3, bytesForNative)
	c.pushUint64(4, bytesForNative)
	c.pushUint64(5, bytesForNative)

	result := MakeSliceOf(allowableTypes, c.getRawBytes())

	expected := []any{
		&interfaceDemoA{IntField: -100, Float64Field: 123.123},
		&interfaceDemoB{},
		&interfaceDemoC{},
		&interfaceDemoD{},
		&interfaceDemoE{},
		&interfaceDemoF{},
	}

	assert.Equal(t, expected, result)
}
