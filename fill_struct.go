package fuzzhelper

import (
	"fmt"
	"reflect"
	"slices"
)

type fillFunc func() []fillFunc

func newFillFunc(value reflect.Value, c *ByteConsumer, tags fuzzTags) fillFunc {
	return func() []fillFunc {
		return fill(value, c, tags)
	}
}

func Fill(root any, c *ByteConsumer) {
	fillFuncs := fill(reflect.ValueOf(root), c, newEmptyFuzzTags())

	values := newDequeue[fillFunc]()
	slices.Reverse(fillFuncs)
	values.addMany(fillFuncs)

	for values.len() != 0 {
		ff := values.popBack()
		fillFuncs := ff()
		slices.Reverse(fillFuncs)
		values.addMany(fillFuncs)
	}

	println("")
}

func fill(value reflect.Value, c *ByteConsumer, tags fuzzTags) []fillFunc {
	if c.Len() == 0 {
		// There are no more bytes to use to fill data
		return []fillFunc{}
	}

	switch value.Kind() {
	case reflect.Bool:
		fillBool(value, c)
		return []fillFunc{}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fillInt(value, c, tags)
		return []fillFunc{}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fillUint(value, c, tags)
		return []fillFunc{}

	case reflect.Uintptr:
		// Uintptr is ignored
		return []fillFunc{}

	case reflect.Float32, reflect.Float64:
		fillFloat(value, c, tags)
		return []fillFunc{}

	case reflect.Complex64, reflect.Complex128:
		// Complex are ignored
		// Only because I don't use them, and I don't think many people use them
		// If needed support should be easy to add
		return []fillFunc{}

	case reflect.Array:
		return fillArray(value, c)

	case reflect.Chan:
		return fillChan(value, c, tags)

	case reflect.Func:
		// functions are ignored
		return []fillFunc{}

	case reflect.Interface:
		// Can't do anything here - we can't instantiate an interface type
		// We don't know which type to create here
		return []fillFunc{}

	case reflect.Map:
		return fillMap(value, c, tags)

	case reflect.Pointer:
		return fillPointer(value, c)

	case reflect.Slice:
		return fillSlice(value, c, tags)

	case reflect.String:
		fillString(value, c, tags)
		return []fillFunc{}

	case reflect.Struct:
		return fillStruct(value, c)

	case reflect.UnsafePointer:
		// Unsafe pointers are just ignored
		return []fillFunc{}

	default:
		fmt.Printf("Unsupported kind %s\n", value.Kind())
	}

	panic("unreachable")
}

func canSet(value reflect.Value) bool {
	// The initial value passed into Fill method must be an
	// instantiated struct/map/array/slice or a pointer to one of
	// these.  Once we drill past this unsettable level we will
	// fill in values recursively as we find them.
	if value.CanSet() {
		println(": can set")
		return true
	}

	println(": can't set")
	return false
}

func fillString(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		val := c.Uint64(BytesForNative)
		str := tags.stringValues[val%uint64(len(tags.stringValues))]

		print("string ", len(str))
		if !canSet(value) {
			return
		}

		value.SetString(str)
		return
	}

	lengthVal := int(c.Int64(BytesForNative))
	strLength := tags.fitStringLength(lengthVal)

	print("string ", strLength)
	if !canSet(value) {
		return
	}

	val := c.String(strLength)
	value.SetString(val)
}

func fillBool(value reflect.Value, c *ByteConsumer) {
	print("bool")
	if !canSet(value) {
		return
	}
	val := c.Bool()
	value.SetBool(val)
}

func fillInt(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	print("int")

	// First check there is a list of valid int values
	if len(tags.intValues) != 0 {
		val := c.Uint64(BytesForNative)
		intVal := tags.intValues[val%uint64(len(tags.intValues))]
		if !canSet(value) {
			return
		}

		value.SetInt(intVal)
		return
	}

	if !canSet(value) {
		return
	}
	val := c.Int64(value.Type().Size())
	fittedVal := tags.fitIntVal(val)
	value.SetInt(fittedVal)
}

func fillUint(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	print("uint")

	// First check there is a list of valid uint values
	if len(tags.uintValues) != 0 {
		val := c.Uint64(BytesForNative)
		uintVal := tags.uintValues[val%uint64(len(tags.uintValues))]
		if !canSet(value) {
			return
		}

		value.SetUint(uintVal)
		return
	}

	if !canSet(value) {
		return
	}
	val := c.Uint64(value.Type().Size())
	fittedVal := tags.fitUintVal(val)
	value.SetUint(fittedVal)
}

func fillFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	print("float")

	// First check there is a list of valid uint values
	if len(tags.floatValues) != 0 {
		val := c.Uint64(BytesForNative)
		floatVal := tags.floatValues[val%uint64(len(tags.floatValues))]
		if !canSet(value) {
			return
		}

		value.SetFloat(floatVal)
		return
	}

	if !canSet(value) {
		return
	}

	val := c.Float64(value.Type().Size())
	fittedVal := tags.fitFloatVal(val)
	value.SetFloat(fittedVal)
}

func fillStruct(value reflect.Value, c *ByteConsumer) []fillFunc {
	print("struct")
	canSet(value)

	newValues := []fillFunc{}
	vType := value.Type()
	for i := 0; i < vType.NumField(); i++ {
		vField := value.Field(i)
		tField := vType.Field(i)
		tags := newFuzzTags(value, tField)
		newValues = append(newValues, newFillFunc(vField, c, tags))
	}

	return newValues
}

func fillPointer(value reflect.Value, c *ByteConsumer) []fillFunc {
	print("pointer")
	if !canSet(value) && value.IsNil() {
		return []fillFunc{}
	}

	if value.IsNil() {
		// If the value is nil - allocate a value for it to point to
		pType := value.Type()
		vType := pType.Elem()
		newVal := reflect.New(vType)
		value.Set(newVal)
	}
	return []fillFunc{
		newFillFunc(value.Elem(), c, newEmptyFuzzTags()),
	}
}

func fillSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags) []fillFunc {
	val := int(c.Int64(BytesForNative))
	sliceLen := tags.fitSliceLengthVal(val)

	print("slice ", sliceLen)
	if !canSet(value) && value.IsNil() {
		return []fillFunc{}
	}

	if value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), sliceLen, sliceLen)
		value.Set(newSlice)
	}

	newValues := []fillFunc{}
	for i := 0; i < value.Len(); i++ {
		newValues = append(newValues, newFillFunc(value.Index(i), c, tags))
	}

	return newValues
}

func fillArray(value reflect.Value, c *ByteConsumer) []fillFunc {
	print("array")
	canSet(value)

	newValues := []fillFunc{}
	for i := 0; i < value.Len(); i++ {
		newValues = append(newValues, newFillFunc(value.Index(i), c, newEmptyFuzzTags()))
	}
	return newValues
}

// TODO there is a bug here where if the map cannot be set but is non-nil this function will try to set it
func fillMap(value reflect.Value, c *ByteConsumer, tags fuzzTags) []fillFunc {
	val := int(c.Int64(BytesForNative))
	mapLen := tags.fitMapLength(val)

	print("map ", mapLen)
	if !canSet(value) && value.IsNil() {
		return []fillFunc{}
	}

	mapType := value.Type()
	keyType := mapType.Key()
	valType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)
	newValues := []fillFunc{}

	for range mapLen {
		// Create the key
		mapKeyP := reflect.New(keyType)
		mapKey := mapKeyP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the key
		newValues = append(newValues, newFillFunc(mapKey, c, tags))

		// Create the value
		mapValP := reflect.New(valType)
		mapVal := mapValP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the value
		newValues = append(newValues, newFillFunc(mapVal, c, tags))

		// Add key/val to map
		newValues = append(newValues, func() []fillFunc {
			newMap.SetMapIndex(mapKey, mapVal)
			return []fillFunc{}
		})
	}

	value.Set(newMap)

	return newValues
}

// TODO there is a bug here, if the channel can't be set, but is non-nil we will still try to set it
func fillChan(value reflect.Value, c *ByteConsumer, tags fuzzTags) []fillFunc {
	val := int(c.Int64(BytesForNative))
	chanLen := tags.fitChanLength(val)

	print("chan ", chanLen)
	if !canSet(value) && value.IsNil() {
		return []fillFunc{}
	}

	chanType := value.Type()
	valType := chanType.Elem()

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), chanLen)
	newValues := []fillFunc{}

	for range chanLen {
		// Create an element for that channel
		newValP := reflect.New(valType)
		newVal := newValP.Elem()
		// Note here that the tags used to create this chan are also
		// used to create the values added to the channel
		newValues = append(newValues, newFillFunc(newVal, c, tags))

		// Put the element on the channel
		newValues = append(newValues, func() []fillFunc {
			newChan.Send(newVal)
			return []fillFunc{}
		})
	}

	value.Set(newChan)

	return newValues
}
