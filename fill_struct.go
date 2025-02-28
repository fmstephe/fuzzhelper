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
	if !value.CanSet() {
		println("can't set")
		// The initial value passed into this method must be an
		// instantiated struct/map/array/slice or a pointer to one of
		// these.  Once we drill past this unsettable level we will
		// fill in values recursively as we find them.
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			println("slice/array")
			if value.IsNil() {
				return
			}
			fillSliceArray(value)
			return

		case reflect.Map:
			// TODO handle maps
			return

		case reflect.Struct:
			println("struct")
			fillStruct(value)
			return

		case reflect.Pointer, reflect.Interface:
			if value.IsNil() {
				// If the value is unsettable and nil, there's nothing we can do
				return
			}
			// Not nil, see if we can set anything after following the pointer/interface
			println("pointer")
			fill(value.Elem())
		}
	} else {
		println("can set")
		// Switch on the next expected kind
		switch value.Kind() {
		case reflect.String:
			println("string")
			value.SetString("string")

		case reflect.Bool:
			println("bool")
			value.SetBool(true)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			println("int")
			value.SetInt(-1)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			println("uint")
			value.SetUint(1)

		case reflect.Float32, reflect.Float64:
			println("float")
			value.SetFloat(1.234)

		case reflect.Complex64, reflect.Complex128:
			println("complex")
			value.SetComplex(1 + 2i)

		case reflect.Array, reflect.Slice:
			fillSliceArray(value)
			return

		case reflect.Map:
			// TODO handle maps
			return

		case reflect.Struct:
			println("struct")
			fillStruct(value)

		case reflect.Chan:
			// TODO handle channels

		case reflect.Pointer:
			// 1: Create an instance of this pointer type
			// 2: Follow that pointer and try to fill it
			println("pointer")
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
}

func fillStruct(value reflect.Value) {
	vType := value.Type()
	for i := 0; i < vType.NumField(); i++ {
		//tField := vType.Field(i)
		// TODO do some checking here on the field's tags
		vField := value.Field(i)
		fill(vField)
	}
}

func fillPointer(value reflect.Value) {
	if value.IsNil() {
		// If the value is nil - allocate a value for it to point to
		pType := value.Type()
		vType := pType.Elem()
		newVal := reflect.New(vType)
		value.Set(newVal)
	}
	fill(value.Elem())
}

func fillSliceArray(value reflect.Value) {
	if value.Kind() == reflect.Slice && value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), 4, 4)
		value.Set(newSlice)
	}
	for i := 0; i < value.Len(); i++ {
		fill(value.Index(i))
	}
}
