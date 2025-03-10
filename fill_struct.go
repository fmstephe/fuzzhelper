package fuzzhelper

import (
	"fmt"
	"reflect"
)

func Fill(value any, c *ByteConsumer) {
	fill(reflect.ValueOf(value), c, newEmptyFuzzTags())
	println("")
}

func fill(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	if c.Len() == 0 {
		// There are no more bytes to use to fill data
		return
	}

	switch value.Kind() {
	case reflect.Bool:
		fillBool(value, c)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fillInt(value, c, tags)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fillUint(value, c, tags)

	case reflect.Uintptr:
		// Uintptr is ignored

	case reflect.Float32, reflect.Float64:
		fillFloat(value, c)

	case reflect.Complex64, reflect.Complex128:
		fillComplex(value, c)

	case reflect.Array:
		fillArray(value, c)

	case reflect.Chan:
		fillChan(value, c, tags)

	case reflect.Func:
		// functions are ignored

	case reflect.Interface:
		// Can't do anything here - we can't instantiate an interface type
		// We don't know which type to create here

	case reflect.Map:
		fillMap(value, c, tags)

	case reflect.Pointer:
		fillPointer(value, c)

	case reflect.Slice:
		fillSlice(value, c, tags)

	case reflect.String:
		fillString(value, c, tags)

	case reflect.Struct:
		fillStruct(value, c)

	case reflect.UnsafePointer:
		// Unsafe pointers are just ignored

	default:
		fmt.Printf("Unsupported kind %s\n", value.Kind())
	}
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

func fillFloat(value reflect.Value, c *ByteConsumer) {
	print("float")
	if !canSet(value) {
		return
	}
	val := c.Float64(value.Type().Size())
	value.SetFloat(val)
}

func fillComplex(value reflect.Value, c *ByteConsumer) {
	print("complex")
	if !canSet(value) {
		return
	}
	value.SetComplex(1 + 2i)
}

func fillStruct(value reflect.Value, c *ByteConsumer) {
	print("struct")
	canSet(value)

	vType := value.Type()
	for i := 0; i < vType.NumField(); i++ {
		vField := value.Field(i)
		tField := vType.Field(i)
		tags := newFuzzTags(value, tField)
		fill(vField, c, tags)
	}
}

func fillPointer(value reflect.Value, c *ByteConsumer) {
	print("pointer")
	if !canSet(value) && value.IsNil() {
		return
	}

	if value.IsNil() {
		// If the value is nil - allocate a value for it to point to
		pType := value.Type()
		vType := pType.Elem()
		newVal := reflect.New(vType)
		value.Set(newVal)
	}
	fill(value.Elem(), c, newEmptyFuzzTags())
}

func fillSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	val := int(c.Int64(BytesForNative))
	sliceLen := tags.fitSliceLengthVal(val)

	print("slice ", sliceLen)
	if !canSet(value) && value.IsNil() {
		return
	}

	if value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), sliceLen, sliceLen)
		value.Set(newSlice)
	}

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i), c, newEmptyFuzzTags())
	}
}

func fillArray(value reflect.Value, c *ByteConsumer) {
	print("array")
	canSet(value)

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i), c, newEmptyFuzzTags())
	}
}

func fillMap(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	val := int(c.Int64(BytesForNative))
	mapLen := tags.fitMapLength(val)

	print("map ", mapLen)
	if !canSet(value) && value.IsNil() {
		return
	}

	mapType := value.Type()
	keyType := mapType.Key()
	valType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)

	for range mapLen {
		// Create the key
		mapKeyP := reflect.New(keyType)
		mapKey := mapKeyP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the key
		fill(mapKey, c, tags)

		// Create the value
		mapValP := reflect.New(valType)
		mapVal := mapValP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the value
		fill(mapVal, c, tags)

		// Add key/val to map
		newMap.SetMapIndex(mapKey, mapVal)
	}

	// Set value to be the new map
	value.Set(newMap)
}

func fillChan(value reflect.Value, c *ByteConsumer, tags fuzzTags) {
	val := int(c.Int64(BytesForNative))
	chanLen := tags.fitChanLength(val)

	print("chan ", chanLen)
	if !canSet(value) && value.IsNil() {
		return
	}

	chanType := value.Type()
	valType := chanType.Elem()

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), chanLen)

	for range chanLen {
		// Create an element for that channel
		newValP := reflect.New(valType)
		newVal := newValP.Elem()
		// Note here that the tags used to create this chan are also
		// used to create the values added to the channel
		fill(newVal, c, tags)

		// Put the element on the channel
		newChan.Send(newVal)
	}

	// Set value to be the new channel
	value.Set(newChan)
}
