package fuzzhelper

import (
	"fmt"
	"reflect"
)

type visitFunc func() []visitFunc

type visitCallback interface {
	visitBool(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitInt(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitUint(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitUintptr(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitFloat(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitChan(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitMap(reflect.Value, *ByteConsumer, fuzzTags) int
	visitPointer(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitSlice(reflect.Value, *ByteConsumer, fuzzTags) int
	visitString(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
	visitStruct(reflect.Value, *ByteConsumer, fuzzTags) []visitFunc
}

func newVisitFunc(callback visitCallback, value reflect.Value, c *ByteConsumer, tags fuzzTags) visitFunc {
	return func() []visitFunc {
		//println(fmt.Sprintf("before %#v\n", value.Interface()))
		ffs := visitValue(callback, value, c, tags)
		//println(fmt.Sprintf("after %#v\n", value.Interface()))
		return ffs
	}
}

func visitRoot(callback visitCallback, root any, c *ByteConsumer) {
	visitFuncs := visitValue(callback, reflect.ValueOf(root), c, newEmptyFuzzTags())

	values := newDequeue[visitFunc]()
	values.addMany(visitFuncs)

	for values.len() != 0 {
		ff := values.popFirst()
		visitFuncs := ff()
		values.addMany(visitFuncs)
	}

	//println("")
}

func visitValue(callback visitCallback, value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	if c.Len() == 0 {
		// There are no more bytes to use to visit data
		return []visitFunc{}
	}

	switch value.Kind() {
	case reflect.Bool:
		return callback.visitBool(value, c, tags)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return callback.visitInt(value, c, tags)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return callback.visitUint(value, c, tags)

	case reflect.Uintptr:
		return callback.visitUintptr(value, c, tags)

	case reflect.Float32, reflect.Float64:
		return callback.visitFloat(value, c, tags)

	case reflect.Complex64, reflect.Complex128:
		// Complex values are ignored Only because I don't use them,
		// and I suspect no one else uses them very often. Can be added
		// in if a need is felt
		//return callback.visitComplex(value, c, tags)
		return []visitFunc{}

	case reflect.Array:
		//print("array")
		canSet(value)

		newValues := []visitFunc{}
		for i := 0; i < value.Len(); i++ {
			newValues = append(newValues, visitValue(callback, value.Index(i), c, newEmptyFuzzTags())...)
		}
		return newValues

	case reflect.Chan:
		return callback.visitChan(value, c, tags)

	case reflect.Func:
		// Ignored
		//return callback.visitFunc(value, c, tags)
		return []visitFunc{}

	case reflect.Interface:
		// Ignored
		//return callback.visitInterface(value, c, tags)
		return []visitFunc{}

	case reflect.Map:
		mapLen := callback.visitMap(value, c, tags)

		mapType := value.Type()
		keyType := mapType.Key()
		valType := mapType.Elem()

		newValues := []visitFunc{}
		for range mapLen {
			// Create the key
			mapKeyP := reflect.New(keyType)
			mapKey := mapKeyP.Elem()
			// Note here that the tags used to create this map are also
			// used to create the key
			newValues = append(newValues, visitValue(callback, mapKey, c, tags)...)

			// Create the value
			mapValP := reflect.New(valType)
			mapVal := mapValP.Elem()
			// Note here that the tags used to create this map are also
			// used to create the value
			newValues = append(newValues, visitValue(callback, mapVal, c, tags)...)

			// Add key/val to map
			//println("setting map element")
			value.SetMapIndex(mapKey, mapVal)
		}

		return newValues

	case reflect.Pointer:
		return callback.visitPointer(value, c, tags)

	case reflect.Slice:
		sliceLen := callback.visitSlice(value, c, tags)

		newValues := []visitFunc{}
		for i := range sliceLen {
			newValues = append(newValues, visitValue(callback, value.Index(i), c, tags)...)
		}
		return newValues

	case reflect.String:
		return callback.visitString(value, c, tags)

	case reflect.Struct:
		return callback.visitStruct(value, c, tags)

	case reflect.UnsafePointer:
		// Ignored
		//return callback.visitUnsafePointer(value, c, tags)
		return []visitFunc{}

	default:
		panic(fmt.Errorf("Unsupported kind %s\n", value.Kind()))
	}

	panic("unreachable")
}
