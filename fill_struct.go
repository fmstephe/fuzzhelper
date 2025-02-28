package fuzzhelper

import (
	"fmt"
	"reflect"
)

func Fill(value any) {
	fill(reflect.ValueOf(value))
	println("")
}

func fill(value reflect.Value) {
	// Switch on the next expected kind
	switch value.Kind() {
	case reflect.String:
		fillString(value)

	case reflect.Bool:
		fillBool(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fillInt(value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fillUint(value)

	case reflect.Float32, reflect.Float64:
		fillFloat(value)

	case reflect.Complex64, reflect.Complex128:
		fillComplex(value)

	case reflect.Slice:
		fillSlice(value)
		return

	case reflect.Array:
		fillArray(value)
		return

	case reflect.Map:
		fillMap(value)
		return

	case reflect.Struct:
		fillStruct(value)

	case reflect.Chan:
		fillChan(value)

	case reflect.Pointer:
		fillPointer(value)

	case reflect.Interface:
		// Can't do anything here - we can't instantiate an interface type
		// We don't know which type to create here

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

func fillString(value reflect.Value) {
	print("string")
	if !canSet(value) {
		return
	}
	value.SetString("string")
}

func fillBool(value reflect.Value) {
	print("bool")
	if !canSet(value) {
		return
	}
	value.SetBool(true)
}

func fillInt(value reflect.Value) {
	print("int")
	if !canSet(value) {
		return
	}
	value.SetInt(-1)
}

func fillUint(value reflect.Value) {
	print("uint")
	if !canSet(value) {
		return
	}
	value.SetUint(1)
}

func fillFloat(value reflect.Value) {
	print("float")
	if !canSet(value) {
		return
	}
	value.SetFloat(1.234)
}

func fillComplex(value reflect.Value) {
	print("complex")
	if !canSet(value) {
		return
	}
	value.SetComplex(1 + 2i)
}

func fillStruct(value reflect.Value) {
	print("struct")
	canSet(value)

	vType := value.Type()
	for i := 0; i < vType.NumField(); i++ {
		//tField := vType.Field(i)
		// TODO do some checking here on the field's tags
		vField := value.Field(i)
		fill(vField)
	}
}

func fillPointer(value reflect.Value) {
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
	fill(value.Elem())
}

func fillSlice(value reflect.Value) {
	print("slice")
	if !canSet(value) && value.IsNil() {
		return
	}

	if value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), 4, 4)
		value.Set(newSlice)
	}

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i))
	}
}

func fillArray(value reflect.Value) {
	print("array")
	canSet(value)

	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i))
	}
}

func fillMap(value reflect.Value) {
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
	fill(mapKey)

	// Create the value
	mapValP := reflect.New(valType)
	mapVal := mapValP.Elem()
	fill(mapVal)

	// Add key/val to map
	newMap.SetMapIndex(mapKey, mapVal)

	// Set value to be the new map
	value.Set(newMap)
}

func fillChan(value reflect.Value) {
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
	fill(newVal)

	// Put the element on the channel
	newChan.Send(newVal)

	// Set value to be the new channel
	value.Set(newChan)
}
