package fuzzhelper

import (
	"reflect"
	"strconv"
	"strings"
)

const defaultLengthMin = 0
const defaultLengthMax = 20

type fuzzTags struct {
	// Debugging field containing the fieldName of the struct field the tag was
	// taken from
	fieldName string
	//
	intRange intTagRange
	//
	uintRange uintTagRange
	//
	floatRange floatTagRange
	//
	sliceRange lengthTagRange
	//
	stringRange lengthTagRange
	//
	mapRange lengthTagRange
	//
	intValuesMethod string
	intValues       []int64
	//
	uintValuesMethod string
	uintValues       []uint64
	//
	floatValuesMethod string
	floatValues       []float64
	//
	stringValuesMethod string
	stringValues       []string
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

	if intValues, methodName, ok := callMethodFromTag[[]int64](structVal, field, "fuzz-int-method"); ok {
		t.intValuesMethod = methodName
		t.intValues = intValues
	}

	if uintValues, methodName, ok := callMethodFromTag[[]uint64](structVal, field, "fuzz-uint-method"); ok {
		t.uintValuesMethod = methodName
		t.uintValues = uintValues
	}

	if floatValues, methodName, ok := callMethodFromTag[[]float64](structVal, field, "fuzz-float-method"); ok {
		t.floatValuesMethod = methodName
		t.floatValues = floatValues
	}

	if stringValues, methodName, ok := callMethodFromTag[[]string](structVal, field, "fuzz-string-method"); ok {
		t.stringValuesMethod = methodName
		t.stringValues = stringValues
	}

	return t
}

func newEmptyFuzzTags() fuzzTags {
	return fuzzTags{}
}

func getFloat64MinMax(field reflect.StructField, tag string) (minVal, maxVal float64, found bool) {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return 0, 0, false
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
		return 0, 0, false
	}

	minVal, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return 0, 0, false
	}

	maxVal, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return 0, 0, false
	}

	//println("float64 min max", tag, minVal, maxVal)
	return minVal, maxVal, true
}

func getInt64MinMax(field reflect.StructField, tag string) (minVal, maxVal int64, found bool) {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return 0, 0, false
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
		return 0, 0, false
	}

	minVal, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return 0, 0, false
	}

	maxVal, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return 0, 0, false
	}

	//println("int64 min max", tag, minVal, maxVal)
	return minVal, maxVal, true
}

func getUint64MinMax(field reflect.StructField, tag string) (minVal, maxVal uint64, found bool) {
	//println(field.Tag)

	valStr, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return 0, 0, false
	}

	parts := strings.Split(valStr, ",")
	if len(parts) != 2 {
		//println("bad min max tag", valStr)
	}

	minVal, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		//println("bad min tag value", valStr)
		return 0, 0, false
	}

	maxVal, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		//println("bad max tag value", valStr)
		return 0, 0, false
	}

	//println("uint64 min max", tag, minVal, maxVal)
	return minVal, maxVal, true
}

func callMethodFromTag[T any](structVal reflect.Value, field reflect.StructField, tag string) (val T, methodName string, found bool) {

	methodName, ok := field.Tag.Lookup(tag)
	if !ok {
		//println("no tag found: ", tag, field.Name)
		return val, methodName, false
	}

	if !isExported(methodName) {
		//println("method is not exported, can't be called: ", methodName, field.Name, structVal.Type().String())
		return val, methodName, false
	}
	// Try to get the method from the struct
	// We look for pointer receiver method first, then value receivers
	// We it in this order under the assumption that people usually use pointer receivers
	method := structVal.Addr().MethodByName(methodName)
	if !method.IsValid() {
		method = structVal.MethodByName(methodName)
		if !method.IsValid() {
			//println("no method found: ", methodName, field.Name, structVal.Type().String())
			return val, methodName, false
		}
	}

	methodType := method.Type()
	if methodType.NumIn() != 0 {
		//println(fmt.Sprintf("expected method with no args, method requires %d args", method.Type().NumIn()), methodName, field.Name)
		return val, methodName, false
	}

	if methodType.NumOut() != 1 {
		//println(fmt.Sprintf("expected method returning 1 value, method returns %d value(s)", method.Type().NumOut()), methodName, field.Name)
		return val, methodName, false
	}

	returnType := methodType.Out(0)
	if returnType != reflect.TypeFor[T]() {
		//println(fmt.Sprintf("expected method returning %s, method returns %s", reflect.TypeFor[T](), returnType), methodName, field.Name)
	}

	result := method.Call([]reflect.Value{})

	//println("foo")
	return result[0].Interface().(T), methodName, true
}
