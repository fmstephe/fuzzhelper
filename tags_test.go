package fuzzhelper

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intLimitStruct struct {
	IntField   int   `fuzz-int-range:"-10000,10000"`
	Int64Field int64 `fuzz-int-range:"-1000,1000"`
	Int32Field int32 `fuzz-int-range:"-100,100"`
	Int16Field int16 `fuzz-int-range:"-10,10"`
	Int8Field  int8  `fuzz-int-range:"-1,1"`
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
	UintField   uint   `fuzz-uint-range:"10000,20000"`
	Uint64Field uint64 `fuzz-uint-range:"1000,2000"`
	Uint32Field uint32 `fuzz-uint-range:"100,200"`
	Uint16Field uint16 `fuzz-uint-range:"10,20"`
	Uint8Field  uint8  `fuzz-uint-range:"1,2"`
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

func TestFuzzTags_SliceLength(t *testing.T) {
	type testStruct struct {
		DefaultSlice []int
		OneSlice     []int `fuzz-slice-range:"1,1"`
		FiveSlice    []int `fuzz-slice-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})
	// Create slice of size 3
	c.pushInt64(3, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)

	// Create a slice of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushInt64(4, BytesForNative)
	c.pushInt64(1, BytesForNative)

	// Create a slice of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushInt64(10, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)

	expected := testStruct{
		DefaultSlice: []int{1, 2, 3},
		OneSlice:     []int{1},
		FiveSlice:    []int{1, 2, 3, 4},
	}

	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, expected, val)
}

func TestFuzzTags_StringLength(t *testing.T) {
	type testStruct struct {
		DefaultString string
		OneString     string `fuzz-string-range:"1,1"`
		FiveString    string `fuzz-string-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})
	// Create slice of size 3

	c.pushUint64(3, BytesForNative)
	c.pushBytes([]byte("abc"))

	// Create a slice of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushUint64(4, BytesForNative)
	c.pushBytes([]byte("a"))

	// Create a slice of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushUint64(10, BytesForNative)
	c.pushBytes([]byte("abcdefgh"))

	expected := testStruct{
		DefaultString: "abc",
		OneString:     "a",
		FiveString:    "abcd",
	}

	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, expected, val)
}

func TestFuzzTags_MapLength(t *testing.T) {
	type testStruct struct {
		DefaultMap map[int]int
		OneMap     map[int]int `fuzz-map-range:"1,1"`
		FiveMap    map[int]int `fuzz-map-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})

	// Create map of size 3
	c.pushUint64(3, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(-1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(-2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(-3, BytesForNative)

	// Create a map of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushUint64(4, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(-1, BytesForNative)

	// Create a map of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushUint64(10, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(-1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(-2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(-3, BytesForNative)
	c.pushInt64(4, BytesForNative)
	c.pushInt64(-4, BytesForNative)

	expected := testStruct{
		DefaultMap: map[int]int{
			1: -1,
			2: -2,
			3: -3,
		},
		OneMap: map[int]int{
			1: -1,
		},
		FiveMap: map[int]int{
			1: -1,
			2: -2,
			3: -3,
			4: -4,
		},
	}

	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, expected, val)
}

func TestFuzzTags_MapLengthKeysAndValues(t *testing.T) {
	type testStruct struct {
		DefaultMap map[string][]int `fuzz-map-range:"0,5" fuzz-string-range:"0,4" fuzz-slice-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})

	// Create map of size 3
	c.pushUint64(3, BytesForNative)

	// First Key/Value
	// String Key of length 3
	c.pushInt64(3, BytesForNative)
	c.pushBytes([]byte("abc"))
	// Slice Value of length 1
	c.pushInt64(1, BytesForNative)
	c.pushInt64(1, BytesForNative)

	// Second Key/Value
	// String Key of length 4
	c.pushInt64(4, BytesForNative)
	c.pushBytes([]byte("abcd"))
	// Slice Value of length 2
	c.pushInt64(8, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)

	// Third Key/Value
	// String Key of length 2
	c.pushInt64(7, BytesForNative)
	c.pushBytes([]byte("ab"))
	// Slice Value of length 5
	c.pushInt64(11, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)
	c.pushInt64(5, BytesForNative)

	expected := testStruct{
		DefaultMap: map[string][]int{
			"abc":  []int{1},
			"abcd": []int{1, 2},
			"ab":   []int{1, 2, 3, 4, 5},
		},
	}

	val := testStruct{}
	Fill(&val, c)
	assert.Equal(t, expected, val)
}

func TestFuzzTags_ChanLength(t *testing.T) {
	type testStruct struct {
		DefaultChan chan int
		OneChan     chan int `fuzz-chan-range:"1,1"`
		FiveChan    chan int `fuzz-chan-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})

	// Create chan of size 3
	c.pushUint64(3, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)

	// Create a chan of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushUint64(4, BytesForNative)
	c.pushInt64(1, BytesForNative)

	// Create a chan of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushUint64(10, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)

	val := testStruct{}
	Fill(&val, c)

	assert.Equal(t, 3, len(val.DefaultChan))
	assert.Equal(t, 1, <-val.DefaultChan)
	assert.Equal(t, 2, <-val.DefaultChan)
	assert.Equal(t, 3, <-val.DefaultChan)

	assert.Equal(t, 1, len(val.OneChan))
	assert.Equal(t, 1, <-val.OneChan)

	assert.Equal(t, 4, len(val.FiveChan))
	assert.Equal(t, 1, <-val.FiveChan)
	assert.Equal(t, 2, <-val.FiveChan)
	assert.Equal(t, 3, <-val.FiveChan)
	assert.Equal(t, 4, <-val.FiveChan)
}

func TestFuzzTags_ChanLengthKeysAndValues(t *testing.T) {
	type testStruct struct {
		StringChan chan string `fuzz-chan-range:"0,5" fuzz-string-range:"0,4"`
		SliceChan  chan []int  `fuzz-chan-range:"0,5" fuzz-slice-range:"0,5"`
	}

	c := NewByteConsumer([]byte{})

	// StringChan chan of size 3
	c.pushUint64(3, BytesForNative)

	// String of length 3
	c.pushInt64(3, BytesForNative)
	c.pushBytes([]byte("abc"))
	// String of length 2
	c.pushInt64(7, BytesForNative)
	c.pushBytes([]byte("ab"))
	// String of length 4
	c.pushInt64(9, BytesForNative)
	c.pushBytes([]byte("abcd"))

	// SliceChan chan of size 3
	c.pushUint64(3, BytesForNative)

	// Slice Value of length 2
	c.pushInt64(8, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)

	// Slice Value of length 5
	c.pushInt64(11, BytesForNative)
	c.pushInt64(1, BytesForNative)
	c.pushInt64(2, BytesForNative)
	c.pushInt64(3, BytesForNative)
	c.pushInt64(4, BytesForNative)
	c.pushInt64(5, BytesForNative)

	// Slice Value of length 1
	c.pushInt64(1, BytesForNative)
	c.pushInt64(1, BytesForNative)

	val := testStruct{}
	Fill(&val, c)

	assert.Equal(t, 3, len(val.StringChan))
	assert.Equal(t, "abc", <-val.StringChan)
	assert.Equal(t, "ab", <-val.StringChan)
	assert.Equal(t, "abcd", <-val.StringChan)

	assert.Equal(t, 3, len(val.SliceChan))
	assert.Equal(t, []int{1, 2}, <-val.SliceChan)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, <-val.SliceChan)
	assert.Equal(t, []int{1}, <-val.SliceChan)
}

func TestFuzzTags_StringValues(t *testing.T) {
	type testStruct struct {
		StringField0 string `fuzz-string-values:"zero,one,two,three,four,five"`
		StringField1 string `fuzz-string-values:"zero,one,two,three,four,five"`
		StringField2 string `fuzz-string-values:"zero,one,two,three,four,five"`
		StringField3 string `fuzz-string-values:"zero,one,two,three,four,five"`
		StringField4 string `fuzz-string-values:"zero,one,two,three,four,five"`
		StringField5 string `fuzz-string-values:"zero,one,two,three,four,five"`
	}

	c := NewByteConsumer([]byte{})
	c.pushUint64(0, BytesForNative)
	c.pushUint64(1, BytesForNative)
	c.pushUint64(2, BytesForNative)
	c.pushUint64(3, BytesForNative)
	c.pushUint64(4, BytesForNative)
	c.pushUint64(5, BytesForNative)

	val := testStruct{}
	Fill(&val, c)

	assert.Equal(t, "zero", val.StringField0)
	assert.Equal(t, "one", val.StringField1)
	assert.Equal(t, "two", val.StringField2)
	assert.Equal(t, "three", val.StringField3)
	assert.Equal(t, "four", val.StringField4)
	assert.Equal(t, "five", val.StringField5)
}
