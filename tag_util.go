package fuzzhelper

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

type intTagRange struct {
	wasSet bool
	intMin int64
	intMax int64
}

func newIntTagRange(field reflect.StructField, tag string) intTagRange {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return intTagRange{}
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
		return intTagRange{}
	}

	minVal, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return intTagRange{}
	}

	maxVal, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return intTagRange{}
	}

	//println("int64 min max", tag, minVal, maxVal)
	return intTagRange{
		wasSet: true,
		intMin: minVal,
		intMax: maxVal,
	}
}

func (r *intTagRange) fit(val int64) int64 {
	if !r.wasSet {
		return val
	}

	spread := (r.intMax - r.intMin) + 1
	if spread <= 0 {
		return val
	}

	fitted := (absInt(val) % spread) + r.intMin
	//println("int val fitted", val, r.intMin, r.intMax, fitted)

	return fitted
}

func absInt(val int64) int64 {
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

type uintTagRange struct {
	wasSet  bool
	uintMin uint64
	uintMax uint64
}

func newUintTagRange(field reflect.StructField, tag string) uintTagRange {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return uintTagRange{}
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
		return uintTagRange{}
	}

	minVal, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return uintTagRange{}
	}

	maxVal, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return uintTagRange{}
	}

	//println("uint64 min max", tag, minVal, maxVal)
	return uintTagRange{
		wasSet:  true,
		uintMin: minVal,
		uintMax: maxVal,
	}
}

func (r *uintTagRange) fit(val uint64) uint64 {
	if !r.wasSet {
		return val
	}

	spread := (r.uintMax - r.uintMin) + 1
	if spread <= 0 {
		return val
	}

	fitted := (val % spread) + r.uintMin
	//pruintln("uint val fitted", val, r.uintMin, r.uintMax, fitted)

	return fitted
}

type floatTagRange struct {
	wasSet   bool
	floatMin float64
	floatMax float64
}

func newFloatTagRange(field reflect.StructField, tag string) floatTagRange {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return floatTagRange{}
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
		return floatTagRange{}
	}

	minVal, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return floatTagRange{}
	}

	maxVal, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return floatTagRange{}
	}

	//println("float64 min max", tag, minVal, maxVal)
	return floatTagRange{
		wasSet:   true,
		floatMin: minVal,
		floatMax: maxVal,
	}
}

func (r *floatTagRange) fit(val float64) float64 {
	if !r.wasSet {
		return val
	}

	spread := (r.floatMax - r.floatMin)
	if spread <= 0 {
		return val
	}

	// If val is not-a-number then just take the mid-point between min and max
	if math.IsNaN(val) {
		return r.floatMin + (spread / 2)
	}

	// If val is positive infinity then take max
	if math.IsInf(val, 1) {
		return r.floatMax
	}

	// If val is negative infinity then take min
	if math.IsInf(val, -1) {
		return r.floatMin
	}

	fitted := math.Mod(math.Abs(val), spread) + r.floatMin
	//println("float val fitted", val, r.floatMin, r.floatMax, fitted)

	return fitted
}
