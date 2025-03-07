package fuzzhelper

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

const defaultLengthMin = 0
const defaultLengthMax = 20

type fuzzTags struct {
	intMax int64
	intMin int64
	//
	uintMax uint64
	uintMin uint64
	//
	sliceLengthMin uint64
	sliceLengthMax uint64
	//
	stringLengthMin uint64
	stringLengthMax uint64
}

func newFuzzTags(field reflect.StructField) fuzzTags {
	t := newEmptyFuzzTags()

	intMin, intMax, ok := getInt64MinMax(field, "fuzz-int-range")
	if ok {
		t.intMin = intMin
		t.intMax = intMax
	}

	uintMin, uintMax, ok := getUint64MinMax(field, "fuzz-uint-range")
	if ok {
		t.uintMin = uintMin
		t.uintMax = uintMax
	}

	sliceLengthMin, sliceLengthMax, ok := getUint64MinMax(field, "fuzz-slice-range")
	if ok {
		t.sliceLengthMin = sliceLengthMin
		t.sliceLengthMax = sliceLengthMax
	} else {
		t.sliceLengthMin = defaultLengthMin
		t.sliceLengthMax = defaultLengthMax
	}

	stringLengthMin, stringLengthMax, ok := getUint64MinMax(field, "fuzz-string-range")
	if ok {
		t.stringLengthMin = stringLengthMin
		t.stringLengthMax = stringLengthMax
	} else {
		t.stringLengthMin = defaultLengthMin
		t.stringLengthMax = defaultLengthMax
	}

	return t
}

func newEmptyFuzzTags() fuzzTags {
	return fuzzTags{
		intMax: math.MaxInt64,
		intMin: math.MinInt64,
	}
}

func (t *fuzzTags) fitIntVal(val int64) int64 {
	return fitIntValInternal(t.intMin, t.intMax, val)
}

func (t *fuzzTags) fitUintVal(val uint64) uint64 {
	return fitUintValInternal(t.uintMin, t.uintMax, val)
}

func (t *fuzzTags) fitSliceLengthVal(val int) int {
	return fitLengthVal(t.sliceLengthMin, t.sliceLengthMax, val)
}

func (t *fuzzTags) fitStringLength(val int) int {
	return fitLengthVal(t.stringLengthMin, t.stringLengthMax, val)
}

func fitLengthVal(lengthMin, lengthMax uint64, val int) int {
	uintLength := uint64(0)

	if val < 0 {
		uintLength = lengthMin
	} else {
		uintLength = fitUintValInternal(lengthMin, lengthMax, uint64(val))
	}

	// Double check that the value fits inside int
	if uintLength > uint64(math.MaxInt) {
		// If you are creating a slice or a string etc. this value will
		// likely allocate more memory than you have. But for pure
		// simplicity we stick to values which fit within the types
		// used here.
		//
		// If you hit this then your length limits are configured wrong.
		return math.MaxInt
	}

	return int(uintLength)
}

func fitIntValInternal(intMin, intMax, val int64) int64 {
	if intMin == 0 && intMax == 0 {
		return val
	}

	spread := (intMax - intMin) + 1
	if spread <= 0 {
		return val
	}

	fitted := (abs(val) % spread) + intMin
	println("int val fitted", val, intMax, intMin, fitted)

	return fitted
}

func fitUintValInternal(uintMin, uintMax, val uint64) uint64 {
	if uintMin == 0 && uintMax == 0 {
		return val
	}

	spread := (uintMax - uintMin) + 1
	if spread <= 0 {
		return val
	}

	fitted := (val % spread) + uintMin
	println("uint val fitted", val, uintMax, uintMin, fitted)

	return fitted
}

func abs(val int64) int64 {
	if val == math.MinInt64 {
		// taking -math.MinInt64 produces math.MinInt64
		// So we need to special case this value
		return math.MaxInt64
	}
	if val < 0 {
		return -val
	}
	return val
}

func getInt64MinMax(field reflect.StructField, tag string) (minVal, maxVal int64, found bool) {
	println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		println("no tag found: ", tag, field.Name)
		return 0, 0, false
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		println("bad min max tag", valStr)
		return 0, 0, false
	}

	minVal, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		println("bad min tag value", valStr)
		return 0, 0, false
	}

	maxVal, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		println("bad max tag value", valStr)
		return 0, 0, false
	}

	println("int64 min max", tag, minVal, maxVal)
	return minVal, maxVal, true
}

func getUint64MinMax(field reflect.StructField, tag string) (minVal, maxVal uint64, found bool) {
	println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		println("no tag found: ", tag, field.Name)
		return 0, 0, false
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		println("bad min max tag", valStr)
	}

	minVal, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		println("bad min tag value", valStr)
		return 0, 0, false
	}

	maxVal, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		println("bad max tag value", valStr)
		return 0, 0, false
	}

	println("uint64 min max", tag, minVal, maxVal)
	return minVal, maxVal, true
}
