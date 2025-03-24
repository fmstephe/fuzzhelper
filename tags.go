package fuzzhelper

import (
	"reflect"
)

const defaultLengthMin = 0
const defaultLengthMax = 20

type fuzzTags struct {
	// Debugging field containing the fieldName of the struct field the tag was
	// taken from
	fieldName string

	intRange    intTagRange
	uintRange   uintTagRange
	floatRange  floatTagRange
	sliceRange  lengthTagRange
	stringRange lengthTagRange
	mapRange    lengthTagRange

	intValues    valueTag[[]int64]
	uintValues   valueTag[[]uint64]
	floatValues  valueTag[[]float64]
	stringValues valueTag[[]string]
}

func newFuzzTags(structVal reflect.Value, field reflect.StructField) fuzzTags {
	t := newEmptyFuzzTags()

	t.fieldName = field.Name

	t.intRange = newIntTagRange(field, "fuzz-int-range")
	t.uintRange = newUintTagRange(field, "fuzz-uint-range")
	t.floatRange = newFloatTagRange(field, "fuzz-float-range")
	t.sliceRange = newLengthTagRange(field, "fuzz-slice-range", defaultLengthMin, defaultLengthMax)
	t.stringRange = newLengthTagRange(field, "fuzz-string-range", defaultLengthMin, defaultLengthMax)
	t.mapRange = newLengthTagRange(field, "fuzz-map-range", defaultLengthMin, defaultLengthMax)

	t.intValues = newValueTag[[]int64](structVal, field, "fuzz-int-method")
	t.uintValues = newValueTag[[]uint64](structVal, field, "fuzz-uint-method")
	t.floatValues = newValueTag[[]float64](structVal, field, "fuzz-float-method")
	t.stringValues = newValueTag[[]string](structVal, field, "fuzz-string-method")

	return t
}

func newEmptyFuzzTags() fuzzTags {
	return fuzzTags{}
}
