package fuzzhelper

import (
	"math"
	"math/rand"
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
		c := newByteConsumer([]byte{})
		c.pushInt64(int64(i), bytesForNative)
		c.pushInt64(int64(i), bytesFor64)
		c.pushInt64(int64(i), bytesFor32)
		c.pushInt64(int64(i), bytesFor16)
		c.pushInt64(int64(i), bytesFor8)

		val := intLimitStruct{}
		Fill(&val, c.getRawBytes())

		assertIntLimits(t, val)
	}
}

func TestFuzzTags_IntLimits_Positive(t *testing.T) {
	c := newByteConsumer([]byte{})
	// Push a maximum int value for each field
	c.pushInt64(math.MaxInt, bytesForNative)
	c.pushInt64(math.MaxInt64, bytesFor64)
	c.pushInt64(math.MaxInt32, bytesFor32)
	c.pushInt64(math.MaxInt16, bytesFor16)
	c.pushInt64(math.MaxInt8, bytesFor8)

	val := intLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertIntLimits(t, val)
}

func TestFuzzTags_IntLimits_Negative(t *testing.T) {
	c := newByteConsumer([]byte{})
	// Push a minimum int value for each field
	c.pushInt64(math.MinInt, bytesForNative)
	c.pushInt64(math.MinInt64, bytesFor64)
	c.pushInt64(math.MinInt32, bytesFor32)
	c.pushInt64(math.MinInt16, bytesFor16)
	c.pushInt64(math.MinInt8, bytesFor8)

	val := intLimitStruct{}
	Fill(&val, c.getRawBytes())

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
		c := newByteConsumer([]byte{})
		c.pushUint64(uint64(i), bytesForNative)
		c.pushUint64(uint64(i), bytesFor64)
		c.pushUint64(uint64(i), bytesFor32)
		c.pushUint64(uint64(i), bytesFor16)
		c.pushUint64(uint64(i), bytesFor8)

		val := uintLimitStruct{}
		Fill(&val, c.getRawBytes())

		assertUintLimits(t, val)
	}
}

func TestFuzzTags_UintLimits_Positive(t *testing.T) {
	c := newByteConsumer([]byte{})
	// Push a maximum uint value for each field
	c.pushUint64(math.MaxUint, bytesForNative)
	c.pushUint64(math.MaxUint64, bytesFor64)
	c.pushUint64(math.MaxUint32, bytesFor32)
	c.pushUint64(math.MaxUint16, bytesFor16)
	c.pushUint64(math.MaxUint8, bytesFor8)

	val := uintLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertUintLimits(t, val)
}

func TestFuzzTags_UintLimits_Zero(t *testing.T) {
	c := newByteConsumer([]byte{})
	// Push a minimum uint value for each field
	c.pushUint64(0, bytesForNative)
	c.pushUint64(0, bytesFor64)
	c.pushUint64(0, bytesFor32)
	c.pushUint64(0, bytesFor16)
	c.pushUint64(0, bytesFor8)

	val := uintLimitStruct{}
	Fill(&val, c.getRawBytes())

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

type floatLimitStruct struct {
	Float64FieldBigLimit  float64 `fuzz-float-range:"1000,2000"`
	Float64FieldTinyLimit float64 `fuzz-float-range:"0.1,0.2"`
	//
	Float32FieldBigLimit  float32 `fuzz-float-range:"100,200"`
	Float32FieldTinyLimit float32 `fuzz-float-range:"0.01,0.02"`
}

func TestFuzzTags_Float_Rand(t *testing.T) {
	// Generate some random floats
	for i := 0; i <= 40_000; i++ {
		c := newByteConsumer([]byte{})
		c.pushFloat64(float64(rand.Float64()), bytesFor64)
		c.pushFloat64(float64(rand.Float64()), bytesFor64)
		c.pushFloat64(float64(rand.Float32()), bytesFor32)
		c.pushFloat64(float64(rand.Float32()), bytesFor32)

		val := floatLimitStruct{}
		Fill(&val, c.getRawBytes())

		assertFloatLimits(t, val)
	}
}

func TestFuzzTags_Float_Max(t *testing.T) {
	c := newByteConsumer([]byte{})
	c.pushFloat64(math.MaxFloat64, bytesFor64)
	c.pushFloat64(math.MaxFloat64, bytesFor64)
	c.pushFloat64(math.MaxFloat32, bytesFor32)
	c.pushFloat64(math.MaxFloat32, bytesFor32)

	val := floatLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertFloatLimits(t, val)
}

func TestFuzzTags_Float_Min(t *testing.T) {
	c := newByteConsumer([]byte{})
	c.pushFloat64(-math.MaxFloat64, bytesFor64)
	c.pushFloat64(-math.MaxFloat64, bytesFor64)
	c.pushFloat64(-math.MaxFloat32, bytesFor32)
	c.pushFloat64(-math.MaxFloat32, bytesFor32)

	val := floatLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertFloatLimits(t, val)
}

func TestFuzzTags_Float_Zero(t *testing.T) {
	c := newByteConsumer([]byte{})
	c.pushFloat64(0, bytesFor64)
	c.pushFloat64(0, bytesFor64)
	c.pushFloat64(0, bytesFor32)
	c.pushFloat64(0, bytesFor32)

	val := floatLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertFloatLimits(t, val)
}

func TestFuzzTags_Float_NaN(t *testing.T) {
	c := newByteConsumer([]byte{})
	c.pushFloat64(math.NaN(), bytesFor64)
	c.pushFloat64(math.NaN(), bytesFor64)
	c.pushFloat64(math.NaN(), bytesFor32)
	c.pushFloat64(math.NaN(), bytesFor32)

	val := floatLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertFloatLimits(t, val)
}

func TestFuzzTags_Float_Inf(t *testing.T) {
	c := newByteConsumer([]byte{})
	c.pushFloat64(math.Inf(1), bytesFor64)
	c.pushFloat64(math.Inf(-1), bytesFor64)
	c.pushFloat64(math.Inf(1), bytesFor32)
	c.pushFloat64(math.Inf(-1), bytesFor32)

	val := floatLimitStruct{}
	Fill(&val, c.getRawBytes())

	assertFloatLimits(t, val)
}

func assertFloatLimits(t *testing.T, val floatLimitStruct) {
	assert.LessOrEqual(t, val.Float64FieldBigLimit, float64(2000))
	assert.GreaterOrEqual(t, val.Float64FieldBigLimit, float64(1000))
	//
	assert.LessOrEqual(t, val.Float64FieldTinyLimit, float64(0.2))
	assert.GreaterOrEqual(t, val.Float64FieldTinyLimit, float64(0.1))
	//
	assert.LessOrEqual(t, val.Float32FieldBigLimit, float32(200))
	assert.GreaterOrEqual(t, val.Float32FieldBigLimit, float32(100))
	//
	assert.LessOrEqual(t, val.Float32FieldTinyLimit, float32(0.02))
	assert.GreaterOrEqual(t, val.Float32FieldTinyLimit, float32(0.01))
}

func TestFuzzTags_SliceLength(t *testing.T) {
	type testStruct struct {
		DefaultSlice []int
		OneSlice     []int `fuzz-slice-range:"1,1"`
		FiveSlice    []int `fuzz-slice-range:"0,5"`
	}

	c := newByteConsumer([]byte{})
	// Create slice of size 3
	c.pushInt64(3, bytesForNative)
	c.pushInt64(1, bytesForNative)
	c.pushInt64(2, bytesForNative)
	c.pushInt64(3, bytesForNative)

	// Create a slice of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushInt64(4, bytesForNative)
	c.pushInt64(1, bytesForNative)

	// Create a slice of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushInt64(10, bytesForNative)
	c.pushInt64(1, bytesForNative)
	c.pushInt64(2, bytesForNative)
	c.pushInt64(3, bytesForNative)
	c.pushInt64(4, bytesForNative)

	expected := testStruct{
		DefaultSlice: []int{1, 2, 3},
		OneSlice:     []int{1},
		FiveSlice:    []int{1, 2, 3, 4},
	}

	val := testStruct{}
	Fill(&val, c.getRawBytes())
	assert.Equal(t, expected, val)
}

func TestFuzzTags_StringLength(t *testing.T) {
	type testStruct struct {
		DefaultString string
		OneString     string `fuzz-string-range:"1,1"`
		FiveString    string `fuzz-string-range:"0,5"`
	}

	c := newByteConsumer([]byte{})
	// Create slice of size 3

	c.pushUint64(3, bytesForNative)
	c.pushBytes([]byte("abc"))

	// Create a slice of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushUint64(4, bytesForNative)
	c.pushBytes([]byte("a"))

	// Create a slice of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushUint64(10, bytesForNative)
	c.pushBytes([]byte("abcdefgh"))

	expected := testStruct{
		DefaultString: "abc",
		OneString:     "a",
		FiveString:    "abcd",
	}

	val := testStruct{}
	Fill(&val, c.getRawBytes())
	assert.Equal(t, expected, val)
}

func TestFuzzTags_MapLength(t *testing.T) {
	type testStruct struct {
		DefaultMap map[int]int
		OneMap     map[int]int `fuzz-map-range:"1,1"`
		FiveMap    map[int]int `fuzz-map-range:"0,5"`
	}

	c := newByteConsumer([]byte{})

	// Create map of size 3
	c.pushUint64(3, bytesForNative)
	// Values for DefaultMap
	c.pushInt64(1, bytesForNative)
	c.pushInt64(-1, bytesForNative)
	c.pushInt64(2, bytesForNative)
	c.pushInt64(-2, bytesForNative)
	c.pushInt64(3, bytesForNative)
	c.pushInt64(-3, bytesForNative)

	// Create a map of size 1, the length value consumed will be 4, but
	// the length min/max forces the size to 1
	c.pushUint64(4, bytesForNative)
	// Values for OneMap
	c.pushInt64(1, bytesForNative)
	c.pushInt64(-1, bytesForNative)

	// Create a map of size 4, the length value consumed will be 10, but
	// because the max length is 5 the fitted value will be 4
	c.pushUint64(10, bytesForNative)

	// Values for FiveMap
	c.pushInt64(1, bytesForNative)
	c.pushInt64(-1, bytesForNative)
	c.pushInt64(2, bytesForNative)
	c.pushInt64(-2, bytesForNative)
	c.pushInt64(3, bytesForNative)
	c.pushInt64(-3, bytesForNative)
	c.pushInt64(4, bytesForNative)
	c.pushInt64(-4, bytesForNative)

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
	Fill(&val, c.getRawBytes())
	assert.Equal(t, expected, val)
}

func TestFuzzTags_MapLengthKeysAndValues(t *testing.T) {
	type testStruct struct {
		DefaultMap map[string][]int `fuzz-map-range:"0,5" fuzz-string-range:"0,4" fuzz-slice-range:"0,5"`
	}

	c := newByteConsumer([]byte{})

	// Create map of size 3
	c.pushUint64(3, bytesForNative)

	// First Key/Value
	// String Key of length 3
	c.pushInt64(3, bytesForNative)
	c.pushBytes([]byte("abc"))
	// Slice Value of length 1
	c.pushInt64(1, bytesForNative)
	c.pushInt64(1, bytesForNative)

	// Second Key/Value
	// String Key of length 4
	c.pushInt64(4, bytesForNative)
	c.pushBytes([]byte("abcd"))
	// Slice Value of length 2
	c.pushInt64(8, bytesForNative)
	c.pushInt64(1, bytesForNative)
	c.pushInt64(2, bytesForNative)

	// Third Key/Value
	// String Key of length 2
	c.pushInt64(7, bytesForNative)
	c.pushBytes([]byte("ab"))
	// Slice Value of length 5
	c.pushInt64(11, bytesForNative)

	// Values "ab" slice
	c.pushInt64(1, bytesForNative)
	c.pushInt64(2, bytesForNative)
	c.pushInt64(3, bytesForNative)
	c.pushInt64(4, bytesForNative)
	c.pushInt64(5, bytesForNative)

	expected := testStruct{
		DefaultMap: map[string][]int{
			"abc":  []int{1},
			"abcd": []int{1, 2},
			"ab":   []int{1, 2, 3, 4, 5},
		},
	}

	val := testStruct{}
	Fill(&val, c.getRawBytes())
	assert.Equal(t, expected, val)
}

type methodStruct struct {
	StringField0 string `fuzz-string-method:"StringOptions"`
	StringField1 string `fuzz-string-method:"StringOptions"`
	StringField2 string `fuzz-string-method:"StringOptions"`
	StringField3 string `fuzz-string-method:"StringOptions"`
	StringField4 string `fuzz-string-method:"StringOptions"`
	StringField5 string `fuzz-string-method:"StringOptions"`
	//
	IntField0 int64 `fuzz-int-method:"IntOptions"`
	IntField1 int64 `fuzz-int-method:"IntOptions"`
	IntField2 int64 `fuzz-int-method:"IntOptions"`
	IntField3 int64 `fuzz-int-method:"IntOptions"`
	IntField4 int64 `fuzz-int-method:"IntOptions"`
	IntField5 int64 `fuzz-int-method:"IntOptions"`
	//
	UintField0 uint64 `fuzz-uint-method:"UintOptions"`
	UintField1 uint64 `fuzz-uint-method:"UintOptions"`
	UintField2 uint64 `fuzz-uint-method:"UintOptions"`
	UintField3 uint64 `fuzz-uint-method:"UintOptions"`
	UintField4 uint64 `fuzz-uint-method:"UintOptions"`
	UintField5 uint64 `fuzz-uint-method:"UintOptions"`
	//
	FloatField0 float64 `fuzz-float-method:"FloatOptions"`
	FloatField1 float64 `fuzz-float-method:"FloatOptions"`
	FloatField2 float64 `fuzz-float-method:"FloatOptions"`
	FloatField3 float64 `fuzz-float-method:"FloatOptions"`
	FloatField4 float64 `fuzz-float-method:"FloatOptions"`
	FloatField5 float64 `fuzz-float-method:"FloatOptions"`
}

// Pointer receiver method
func (s *methodStruct) StringOptions() []string {
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
func (s methodStruct) IntOptions() []int64 {
	return []int64{
		0,
		-1,
		-2,
		-3,
		-4,
		-5,
	}
}

// Pointer receiver method
func (s *methodStruct) UintOptions() []uint64 {
	return []uint64{
		0,
		10,
		20,
		30,
		40,
		50,
	}
}

// Value receiver method
func (s methodStruct) FloatOptions() []float64 {
	return []float64{
		0.0,
		0.1,
		0.2,
		0.3,
		0.4,
		0.5,
	}
}

func TestFuzzTags_MethodValues(t *testing.T) {
	c := newByteConsumer([]byte{})
	// Values for strings
	c.pushUint64(0, bytesForNative)
	c.pushUint64(1, bytesForNative)
	c.pushUint64(2, bytesForNative)
	c.pushUint64(3, bytesForNative)
	c.pushUint64(4, bytesForNative)
	c.pushUint64(5, bytesForNative)
	// Values for ints
	c.pushUint64(0, bytesForNative)
	c.pushUint64(1, bytesForNative)
	c.pushUint64(2, bytesForNative)
	c.pushUint64(3, bytesForNative)
	c.pushUint64(4, bytesForNative)
	c.pushUint64(5, bytesForNative)
	// Values for uints
	c.pushUint64(0, bytesForNative)
	c.pushUint64(1, bytesForNative)
	c.pushUint64(2, bytesForNative)
	c.pushUint64(3, bytesForNative)
	c.pushUint64(4, bytesForNative)
	c.pushUint64(5, bytesForNative)
	// Values for floats
	c.pushUint64(0, bytesForNative)
	c.pushUint64(1, bytesForNative)
	c.pushUint64(2, bytesForNative)
	c.pushUint64(3, bytesForNative)
	c.pushUint64(4, bytesForNative)
	c.pushUint64(5, bytesForNative)

	val := methodStruct{}
	Fill(&val, c.getRawBytes())

	assert.Equal(t, "zero", val.StringField0)
	assert.Equal(t, "one", val.StringField1)
	assert.Equal(t, "two", val.StringField2)
	assert.Equal(t, "three", val.StringField3)
	assert.Equal(t, "four", val.StringField4)
	assert.Equal(t, "five", val.StringField5)

	assert.Equal(t, int64(0), val.IntField0)
	assert.Equal(t, int64(-1), val.IntField1)
	assert.Equal(t, int64(-2), val.IntField2)
	assert.Equal(t, int64(-3), val.IntField3)
	assert.Equal(t, int64(-4), val.IntField4)
	assert.Equal(t, int64(-5), val.IntField5)

	assert.Equal(t, uint64(0), val.UintField0)
	assert.Equal(t, uint64(10), val.UintField1)
	assert.Equal(t, uint64(20), val.UintField2)
	assert.Equal(t, uint64(30), val.UintField3)
	assert.Equal(t, uint64(40), val.UintField4)
	assert.Equal(t, uint64(50), val.UintField5)

	assert.Equal(t, float64(0.0), val.FloatField0)
	assert.Equal(t, float64(0.1), val.FloatField1)
	assert.Equal(t, float64(0.2), val.FloatField2)
	assert.Equal(t, float64(0.3), val.FloatField3)
	assert.Equal(t, float64(0.4), val.FloatField4)
	assert.Equal(t, float64(0.5), val.FloatField5)
}
