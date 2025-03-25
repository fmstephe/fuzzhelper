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
	// Debugging field containing the fieldName of the struct field the tag was
	// taken from
	fieldName string
	//
	intMax int64
	intMin int64
	//
	uintMax uint64
	uintMin uint64
	//
	floatMax float64
	floatMin float64
	//
	sliceLengthMin uint64
	sliceLengthMax uint64
	//
	stringLengthMin uint64
	stringLengthMax uint64
	//
	mapLengthMin uint64
	mapLengthMax uint64
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

	if intMin, intMax, ok := getInt64MinMax(field, "fuzz-int-range"); ok {
		t.intMin = intMin
		t.intMax = intMax
	}

	if uintMin, uintMax, ok := getUint64MinMax(field, "fuzz-uint-range"); ok {
		t.uintMin = uintMin
		t.uintMax = uintMax
	}

	if floatMin, floatMax, ok := getFloat64MinMax(field, "fuzz-float-range"); ok {
		t.floatMin = floatMin
		t.floatMax = floatMax
	}

	if sliceLengthMin, sliceLengthMax, ok := getUint64MinMax(field, "fuzz-slice-range"); ok {
		t.sliceLengthMin = sliceLengthMin
		t.sliceLengthMax = sliceLengthMax
	} else {
		t.sliceLengthMin = defaultLengthMin
		t.sliceLengthMax = defaultLengthMax
	}

	if stringLengthMin, stringLengthMax, ok := getUint64MinMax(field, "fuzz-string-range"); ok {
		t.stringLengthMin = stringLengthMin
		t.stringLengthMax = stringLengthMax
	} else {
		t.stringLengthMin = defaultLengthMin
		t.stringLengthMax = defaultLengthMax
	}

	if mapLengthMin, mapLengthMax, ok := getUint64MinMax(field, "fuzz-map-range"); ok {
		t.mapLengthMin = mapLengthMin
		t.mapLengthMax = mapLengthMax
	} else {
		t.mapLengthMin = defaultLengthMin
		t.mapLengthMax = defaultLengthMax
	}

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

func (t *fuzzTags) fitMapLength(val int) int {
	return fitLengthVal(t.mapLengthMin, t.mapLengthMax, val)
}

func (t *fuzzTags) fitFloatVal(val float64) float64 {
	return fitFloatValInternal(t.floatMin, t.floatMax, val)
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

	fitted := (absInt(val) % spread) + intMin
	//println("int val fitted", val, intMin, intMax, fitted)

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
	//println("uint val fitted", val, uintMin, uintMax, fitted)

	return fitted
}

func fitFloatValInternal(floatMin, floatMax, val float64) float64 {
	if floatMin == 0 && floatMax == 0 {
		return val
	}

	spread := (floatMax - floatMin)
	if spread <= 0 {
		return val
	}

	// If val is not-a-number then just take the mid-point between min and max
	if math.IsNaN(val) {
		return floatMin + (spread / 2)
	}

	// If val is positive infinity then take max
	if math.IsInf(val, 1) {
		return floatMax
	}

	// If val is negative infinity then take min
	if math.IsInf(val, -1) {
		return floatMin
	}

	fitted := math.Mod(math.Abs(val), spread) + floatMin
	//println("float val fitted", val, floatMin, floatMax, fitted)

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
