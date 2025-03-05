package fuzzhelper

import (
	"math"
	"reflect"
	"strconv"
)

type fuzzTags struct {
	intMax int64
	intMin int64
}

func newFuzzTags(field reflect.StructField) fuzzTags {
	t := newEmptyFuzzTags()

	intMax, ok := getInt64Tag(field, "fuzz-int-max")
	if ok {
		t.intMax = intMax
	}

	intMin, ok := getInt64Tag(field, "fuzz-int-min")
	if ok {
		t.intMin = intMin
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
	spread := t.intMax - t.intMin
	if spread <= 0 {
		return val
	}

	fitted := (abs(val) % spread) + t.intMin
	println("int val fitted", val, t.intMax, t.intMin, fitted)

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

func getInt64Tag(field reflect.StructField, tag string) (int64, bool) {
	println(field.Tag)
	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		println("no tag found")
		return 0, false
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		println("bad tag value", valStr)
		return 0, false
	}

	println(tag, val)
	return val, true
}
