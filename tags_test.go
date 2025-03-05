package fuzzhelper

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intLimitStruct struct {
	IntField   int   `fuzz-int-min:"-10000" fuzz-int-max:"10000"`
	Int64Field int64 `fuzz-int-min:"-1000" fuzz-int-max:"1000"`
	Int32Field int32 `fuzz-int-min:"-100" fuzz-int-max:"100"`
	Int16Field int16 `fuzz-int-min:"-10" fuzz-int-max:"10"`
	Int8Field  int8  `fuzz-int-min:"-1" fuzz-int-max:"1"`
}

func TestFuzzTags_IntLimits(t *testing.T) {
	for i := -20_000; i <= 20_000; i++ {
		c := NewByteConsumer([]byte{})
		c.pushInt64(int64(i), BytesForNative)
		c.pushInt64(int64(i), BytesFor64)
		c.pushInt64(int64(i), BytesFor32)
		c.pushInt64(int64(i), BytesFor16)
		c.pushInt64(int64(i), BytesFor8)

		val := intLimitStruct{}
		Fill(&val, c)

		assertIntLimits(t, val)
	}
}

func TestFuzzTags_IntLimits_Positive(t *testing.T) {
	c := NewByteConsumer([]byte{})
	// Push a maximum int value for each field
	c.pushInt64(math.MaxInt, BytesForNative)
	c.pushInt64(math.MaxInt64, BytesFor64)
	c.pushInt64(math.MaxInt32, BytesFor32)
	c.pushInt64(math.MaxInt16, BytesFor16)
	c.pushInt64(math.MaxInt8, BytesFor8)

	val := intLimitStruct{}
	Fill(&val, c)

	assertIntLimits(t, val)
}

func TestFuzzTags_IntLimits_Negative(t *testing.T) {
	c := NewByteConsumer([]byte{})
	// Push a minimum int value for each field
	c.pushInt64(math.MinInt, BytesForNative)
	c.pushInt64(math.MinInt64, BytesFor64)
	c.pushInt64(math.MinInt32, BytesFor32)
	c.pushInt64(math.MinInt16, BytesFor16)
	c.pushInt64(math.MinInt8, BytesFor8)

	val := intLimitStruct{}
	Fill(&val, c)

	assertIntLimits(t, val)
}

func assertIntLimits(t *testing.T, val intLimitStruct) {
	assert.LessOrEqual(t, val.IntField, 10_000)
	assert.GreaterOrEqual(t, val.IntField, -10_000)

	assert.LessOrEqual(t, val.Int64Field, int64(1000))
	assert.GreaterOrEqual(t, val.Int64Field, int64(-1000))

	assert.LessOrEqual(t, val.Int32Field, int32(100))
	assert.GreaterOrEqual(t, val.Int32Field, int32(-100))

	assert.LessOrEqual(t, val.Int16Field, int16(10))
	assert.GreaterOrEqual(t, val.Int16Field, int16(-10))

	assert.LessOrEqual(t, val.Int8Field, int8(1))
	assert.GreaterOrEqual(t, val.Int8Field, int8(-1))
}

type uintLimitStruct struct {
	UintField   uint   `fuzz-uint-min:"10000" fuzz-uint-max:"20000"`
	Uint64Field uint64 `fuzz-uint-min:"1000" fuzz-uint-max:"2000"`
	Uint32Field uint32 `fuzz-uint-min:"100" fuzz-uint-max:"200"`
	Uint16Field uint16 `fuzz-uint-min:"10" fuzz-uint-max:"20"`
	Uint8Field  uint8  `fuzz-uint-min:"1" fuzz-uint-max:"2"`
}

func TestFuzzTags_UintLimits(t *testing.T) {
	for i := 0; i <= 40_000; i++ {
		c := NewByteConsumer([]byte{})
		c.pushUint64(uint64(i), BytesForNative)
		c.pushUint64(uint64(i), BytesFor64)
		c.pushUint64(uint64(i), BytesFor32)
		c.pushUint64(uint64(i), BytesFor16)
		c.pushUint64(uint64(i), BytesFor8)

		val := uintLimitStruct{}
		Fill(&val, c)

		assertUintLimits(t, val)
	}
}

func TestFuzzTags_UintLimits_Positive(t *testing.T) {
	c := NewByteConsumer([]byte{})
	// Push a maximum uint value for each field
	c.pushUint64(math.MaxUint, BytesForNative)
	c.pushUint64(math.MaxUint64, BytesFor64)
	c.pushUint64(math.MaxUint32, BytesFor32)
	c.pushUint64(math.MaxUint16, BytesFor16)
	c.pushUint64(math.MaxUint8, BytesFor8)

	val := uintLimitStruct{}
	Fill(&val, c)

	assertUintLimits(t, val)
}

func TestFuzzTags_UintLimits_Zero(t *testing.T) {
	c := NewByteConsumer([]byte{})
	// Push a minimum uint value for each field
	c.pushUint64(0, BytesForNative)
	c.pushUint64(0, BytesFor64)
	c.pushUint64(0, BytesFor32)
	c.pushUint64(0, BytesFor16)
	c.pushUint64(0, BytesFor8)

	val := uintLimitStruct{}
	Fill(&val, c)

	assertUintLimits(t, val)
}

func assertUintLimits(t *testing.T, val uintLimitStruct) {
	assert.LessOrEqual(t, val.UintField, uint(20_000))
	assert.GreaterOrEqual(t, val.UintField, uint(10_000))

	assert.LessOrEqual(t, val.Uint64Field, uint64(2000))
	assert.GreaterOrEqual(t, val.Uint64Field, uint64(1000))

	assert.LessOrEqual(t, val.Uint32Field, uint32(200))
	assert.GreaterOrEqual(t, val.Uint32Field, uint32(100))

	assert.LessOrEqual(t, val.Uint16Field, uint16(20))
	assert.GreaterOrEqual(t, val.Uint16Field, uint16(10))

	assert.LessOrEqual(t, val.Uint8Field, uint8(2))
	assert.GreaterOrEqual(t, val.Uint8Field, uint8(1))
}
