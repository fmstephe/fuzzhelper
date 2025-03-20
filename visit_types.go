package fuzzhelper

import (
	"fmt"
	"reflect"
)

type visitFunc func() []visitFunc

type valueVisitor interface {
	visitBool(reflect.Value, *ByteConsumer, fuzzTags)
	visitInt(reflect.Value, *ByteConsumer, fuzzTags)
	visitUint(reflect.Value, *ByteConsumer, fuzzTags)
	visitUintptr(reflect.Value, *ByteConsumer, fuzzTags)
	visitFloat(reflect.Value, *ByteConsumer, fuzzTags)
	visitChan(reflect.Value, *ByteConsumer, fuzzTags) int
	visitMap(reflect.Value, *ByteConsumer, fuzzTags) int
	visitPointer(reflect.Value, *ByteConsumer, fuzzTags)
	visitSlice(reflect.Value, *ByteConsumer, fuzzTags) int
	visitString(reflect.Value, *ByteConsumer, fuzzTags)
}

func newVisitFunc(callback valueVisitor, value reflect.Value, c *ByteConsumer, tags fuzzTags) visitFunc {
	return func() []visitFunc {
		//println(fmt.Sprintf("before %#v\n", value.Interface()))
		ffs := visitValue(callback, value, c, tags)
		//println(fmt.Sprintf("after %#v\n", value.Interface()))
		return ffs
	}
}

func visitRoot(callback valueVisitor, root any, c *ByteConsumer) {
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

func visitValue(callback valueVisitor, value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	if c.Len() == 0 {
		// There are no more bytes to use to visit data
		return []visitFunc{}
	}

	switch value.Kind() {
	case reflect.Bool:
		callback.visitBool(value, c, tags)
		return []visitFunc{}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		callback.visitInt(value, c, tags)
		return []visitFunc{}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		callback.visitUint(value, c, tags)
		return []visitFunc{}

	case reflect.Uintptr:
		callback.visitUintptr(value, c, tags)
		return []visitFunc{}

	case reflect.Float32, reflect.Float64:
		callback.visitFloat(value, c, tags)
		return []visitFunc{}

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
		chanLen := callback.visitChan(value, c, tags)
		valType := value.Type().Elem()
		newValues := []visitFunc{}

		for range chanLen {
			// Create an element for that channel
			newValP := reflect.New(valType)
			newVal := newValP.Elem()
			// Note here that the tags used to create this chan are also
			// used to create the values added to the channel
			newValues = append(newValues, visitValue(callback, newVal, c, tags)...)
			// newVal has been constructed, send it
			value.Send(newVal)
		}

		return newValues

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
		callback.visitPointer(value, c, tags)
		return []visitFunc{
			newVisitFunc(callback, value.Elem(), c, newEmptyFuzzTags()),
		}

	case reflect.Slice:
		sliceLen := callback.visitSlice(value, c, tags)

		newValues := []visitFunc{}
		for i := range sliceLen {
			newValues = append(newValues, visitValue(callback, value.Index(i), c, tags)...)
		}
		return newValues

	case reflect.String:
		callback.visitString(value, c, tags)
		return []visitFunc{}

	case reflect.Struct:
		//print("struct ", value.Type().Name())
		canSet(value)

		newValues := []visitFunc{}
		vType := value.Type()
		for i := 0; i < vType.NumField(); i++ {
			vField := value.Field(i)
			tField := vType.Field(i)
			tags := newFuzzTags(value, tField)
			newValues = append(newValues, visitValue(callback, vField, c, tags)...)
		}

		return newValues

	case reflect.UnsafePointer:
		// Ignored
		//return callback.visitUnsafePointer(value, c, tags)
		return []visitFunc{}

	default:
		panic(fmt.Errorf("Unsupported kind %s\n", value.Kind()))
	}

	panic("unreachable")
}
