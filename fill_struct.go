package fuzzhelper

import (
	"fmt"
	"reflect"
)

func Fill(value any, c *ByteConsumer) {
	fill(reflect.ValueOf(value), c)
	println("")
}

func fill(value reflect.Value, c *ByteConsumer) {
	switch value.Kind() {
	case reflect.Bool:
		fillBool(value, c)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fillInt(value, c)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fillUint(value, c)

	case reflect.Uintptr:
		// Uintptr is ignored

	case reflect.Float32, reflect.Float64:
		fillFloat(value, c)

	case reflect.Complex64, reflect.Complex128:
		fillComplex(value, c)

	case reflect.Array:
		fillArray(value, c)

	case reflect.Chan:
		fillChan(value, c)

	case reflect.Func:
		// functions are ignored

	case reflect.Interface:
		// Can't do anything here - we can't instantiate an interface type
		// We don't know which type to create here

	case reflect.Map:
		fillMap(value, c)

	case reflect.Pointer:
		fillPointer(value, c)

	case reflect.Slice:
		fillSlice(value, c)

	case reflect.String:
		fillString(value, c)

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

func fillString(value reflect.Value, c *ByteConsumer) {
	print("string")
	if !canSet(value) {
		return
	}
	value.SetString("string")
}

func fillBool(value reflect.Value, c *ByteConsumer) {
	print("bool")
	if !canSet(value) {
		return
	}
	val := c.Bool()
	value.SetBool(val)
}

func fillInt(value reflect.Value, c *ByteConsumer) {
	print("int")
	if !canSet(value) {
		return
	}
	val := c.Int64(value.Type().Size())
	value.SetInt(val)
}

func fillUint(value reflect.Value, c *ByteConsumer) {
	print("uint")
	if !canSet(value) {
		return
	}
	val := c.Uint64(value.Type().Size())
	value.SetUint(val)
}

func fillFloat(value reflect.Value, c *ByteConsumer) {
	print("float")
	if !canSet(value) {
		return
	}
	value.SetFloat(1.234)
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
		//tField := vType.Field(i)
		// TODO do some checking here on the field's tags
		vField := value.Field(i)
		fill(vField, c)
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
	fill(value.Elem(), c)
}

func fillSlice(value reflect.Value, c *ByteConsumer) {
	print("slice")
	if !canSet(value) && value.IsNil() {
		return
	}

	if value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), 4, 4)
		value.Set(newSlice)
	}

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i), c)
	}
}

func fillArray(value reflect.Value, c *ByteConsumer) {
	print("array")
	canSet(value)

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i), c)
	}
}

func fillMap(value reflect.Value, c *ByteConsumer) {
	print("map")
	if !canSet(value) && value.IsNil() {
		return
	}

	mapType := value.Type()
	keyType := mapType.Key()
	valType := mapType.Elem()

	// Set only a single element in the map
	// This is all we can do right now because we always fill the same value for every type
	newMap := reflect.MakeMap(mapType)

	// Create the key
	mapKeyP := reflect.New(keyType)
	mapKey := mapKeyP.Elem()
	fill(mapKey, c)

	// Create the value
	mapValP := reflect.New(valType)
	mapVal := mapValP.Elem()
	fill(mapVal, c)

	// Add key/val to map
	newMap.SetMapIndex(mapKey, mapVal)

	// Set value to be the new map
	value.Set(newMap)
}

func fillChan(value reflect.Value, c *ByteConsumer) {
	print("chan")
	if !canSet(value) && value.IsNil() {
		return
	}

	chanType := value.Type()
	valType := chanType.Elem()

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), 1)

	// Create an element for that channel
	newValP := reflect.New(valType)
	newVal := newValP.Elem()
	fill(newVal, c)

	// Put the element on the channel
	newChan.Send(newVal)

	// Set value to be the new channel
	value.Set(newChan)
}
