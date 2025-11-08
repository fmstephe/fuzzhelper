package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test is so large that it gets its own file.
//
// When we process options tags we check that the contents of the slice
// returned can be assigned to the field that tag was attached to. However, we
// then store the values in a slice which holds the widest version of that
// type, e.g.
// 	* int16 options are stored in a slice of int64
// 	* float32 options are stored in a slice of float64.
//	* interface options are stored in a slice of any
// The correct type of the value is restored when we assign it to a field. This
// process _feels_ quite error prone and the reflective code which performs the
// type check and widening conversion is challenging to read. So we
// exhaustively test the type supported type conversions here.

type typeConversionStruct struct {
	StringField0 string `fuzz-string-method:"StringOptions"`
	//
	IntField   int   `fuzz-int-method:"IntOptions"`
	IntField8  int8  `fuzz-int-method:"IntOptions8"`
	IntField16 int16 `fuzz-int-method:"IntOptions16"`
	IntField32 int32 `fuzz-int-method:"IntOptions32"`
	IntField64 int64 `fuzz-int-method:"IntOptions64"`
	//
	UintField   uint   `fuzz-uint-method:"UintOptions"`
	UintField8  uint8  `fuzz-uint-method:"UintOptions8"`
	UintField16 uint16 `fuzz-uint-method:"UintOptions16"`
	UintField32 uint32 `fuzz-uint-method:"UintOptions32"`
	UintField64 uint64 `fuzz-uint-method:"UintOptions64"`
	//
	FloatField32 float32 `fuzz-float-method:"FloatOptions32"`
	FloatField64 float64 `fuzz-float-method:"FloatOptions64"`
	//
	InterfaceFieldDemo interfaceDemo `fuzz-interface-method:"InterfaceOptionsDemo"`
	InterfaceFieldAny  any           `fuzz-interface-method:"InterfaceOptionsAny"`
}

// Pointer receiver method
func (s *typeConversionStruct) StringOptions() []string {
	return []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
	}
}

// Value receiver method
func (s typeConversionStruct) IntOptions() []int {
	return []int{
		0,
		-1,
		-2,
		-3,
		-4,
		-5,
	}
}

// Value receiver method
func (s typeConversionStruct) IntOptions8() []int8 {
	return []int8{
		0,
		-8,
		-16,
		-24,
		-32,
		-40,
	}
}

// Value receiver method
func (s typeConversionStruct) IntOptions16() []int16 {
	return []int16{
		0,
		-16,
		-32,
		-48,
		-64,
		-80,
	}
}

// Value receiver method
func (s typeConversionStruct) IntOptions32() []int32 {
	return []int32{
		0,
		-32,
		-64,
		-96,
		-128,
		-160,
	}
}

// Value receiver method
func (s typeConversionStruct) IntOptions64() []int64 {
	return []int64{
		0,
		-64,
		-128,
		-192,
		-256,
		-320,
	}
}

// Value receiver method
func (s typeConversionStruct) UintOptions() []uint {
	return []uint{
		0,
		1,
		2,
		3,
		4,
		5,
	}
}

// Value receiver method
func (s typeConversionStruct) UintOptions8() []uint8 {
	return []uint8{
		0,
		8,
		16,
		24,
		32,
		40,
	}
}

// Value receiver method
func (s typeConversionStruct) UintOptions16() []uint16 {
	return []uint16{
		0,
		16,
		32,
		48,
		64,
		80,
	}
}

// Value receiver method
func (s typeConversionStruct) UintOptions32() []uint32 {
	return []uint32{
		0,
		32,
		64,
		96,
		128,
		160,
	}
}

// Value receiver method
func (s typeConversionStruct) UintOptions64() []uint64 {
	return []uint64{
		0,
		64,
		128,
		192,
		256,
		320,
	}
}

// Value receiver method
func (s typeConversionStruct) FloatOptions32() []float32 {
	return []float32{
		0.0,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
	}
}

// Value receiver method
func (s typeConversionStruct) FloatOptions64() []float64 {
	return []float64{
		0.0,
		-0.1,
		-0.2,
		-0.3,
		-0.4,
		-0.5,
	}
}

// Pointer receiver method
func (s *typeConversionStruct) InterfaceOptionsDemo() []interfaceDemo {
	return []interfaceDemo{
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
}

// Pointer receiver method
func (s *typeConversionStruct) InterfaceOptionsAny() []any {
	return []any{
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
}

func TestFuzzTags_MethodValues_TypeConversions(t *testing.T) {
	c := newByteConsumer([]byte{})
	// string
	c.pushUint64(0, bytesForNative)

	// int
	c.pushUint64(0, bytesForNative)
	// int8
	c.pushUint64(1, bytesForNative)
	// int16
	c.pushUint64(2, bytesForNative)
	// int32
	c.pushUint64(3, bytesForNative)
	// int64
	c.pushUint64(4, bytesForNative)

	// uint
	c.pushUint64(0, bytesForNative)
	// uint8
	c.pushUint64(1, bytesForNative)
	// uint16
	c.pushUint64(2, bytesForNative)
	// uint32
	c.pushUint64(3, bytesForNative)
	// uint64
	c.pushUint64(4, bytesForNative)

	// float32
	c.pushUint64(0, bytesForNative)
	// float64
	c.pushUint64(1, bytesForNative)

	// InterfaceFieldDemo
	c.pushUint64(0, bytesForNative)
	// InterfaceFieldAny
	c.pushUint64(1, bytesForNative)

	// Values for InterfaceFieldDemo (which was chosen as interfaceDemoA)
	c.pushInt64(-100, bytesForNative)
	c.pushFloat64(123.123, bytesForNative)

	val := typeConversionStruct{}
	Fill(&val, c.getRawBytes())

	assert.Equal(t, "zero", val.StringField0)

	assert.Equal(t, int(0), val.IntField)
	assert.Equal(t, int8(-8), val.IntField8)
	assert.Equal(t, int16(-32), val.IntField16)
	assert.Equal(t, int32(-96), val.IntField32)
	assert.Equal(t, int64(-256), val.IntField64)

	assert.Equal(t, uint(0), val.UintField)
	assert.Equal(t, uint8(8), val.UintField8)
	assert.Equal(t, uint16(32), val.UintField16)
	assert.Equal(t, uint32(96), val.UintField32)
	assert.Equal(t, uint64(256), val.UintField64)

	assert.Equal(t, float32(0.0), val.FloatField32)
	assert.Equal(t, float64(-0.1), val.FloatField64)

	assert.Equal(t, "interfaceDemoA", val.InterfaceFieldDemo.InterfaceMethod())
	assert.Equal(t, "interfaceDemoB", val.InterfaceFieldAny.(interfaceDemo).InterfaceMethod())
}
