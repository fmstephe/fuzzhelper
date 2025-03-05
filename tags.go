package fuzzhelper

import (
	"math"
	"reflect"
	"strconv"
)

type fuzzTags struct {
	intMax  int64
	intMin  int64
	uintMax uint64
	uintMin uint64
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

	uintMax, ok := getUint64Tag(field, "fuzz-uint-max")
	if ok {
		t.uintMax = uintMax
	}

	uintMin, ok := getUint64Tag(field, "fuzz-uint-min")
	if ok {
		t.uintMin = uintMin
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

func (t *fuzzTags) fitUintVal(val uint64) uint64 {
	spread := t.uintMax - t.uintMin
	if spread <= 0 {
		return val
	}

	fitted := (val % spread) + t.uintMin
	println("uint val fitted", val, t.uintMax, t.uintMin, fitted)

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
		println("no tag found: ", tag)
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

func getUint64Tag(field reflect.StructField, tag string) (uint64, bool) {
	println(field.Tag)
	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		println("no tag found: ", tag)
		return 0, false
	}

	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		println("bad tag value", valStr)
		return 0, false
	}

	println(tag, val)
	return val, true
}
