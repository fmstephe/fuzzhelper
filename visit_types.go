package fuzzhelper

import (
	"fmt"
	"reflect"
)

type visitFunc func() []visitFunc

type valueVisitor interface {
	canGrowRootSlice() bool
	visitBool(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitInt(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitUint(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitUintptr(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitFloat(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitComplex(reflect.Value, fuzzTags, valuePath)
	visitArray(reflect.Value, fuzzTags, valuePath)
	visitChan(reflect.Value, fuzzTags, valuePath)
	visitFunc(reflect.Value, fuzzTags, valuePath)
	visitInterface(reflect.Value, fuzzTags, valuePath)
	visitMap(reflect.Value, *byteConsumer, fuzzTags, valuePath) int
	visitPointer(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitSlice(reflect.Value, *byteConsumer, fuzzTags, valuePath) int
	visitString(reflect.Value, *byteConsumer, fuzzTags, valuePath)
	visitStruct(reflect.Value, fuzzTags, valuePath) bool
	visitUnsafePointer(reflect.Value, fuzzTags, valuePath)
}

func newVisitFunc(callback valueVisitor, value reflect.Value, c *byteConsumer, tags fuzzTags, path valuePath) visitFunc {
	return func() []visitFunc {
		//println(fmt.Sprintf("before %#v\n", value.Interface()))
		ffs := visitValue(callback, value, c, tags, path)
		//println(fmt.Sprintf("after %#v\n", value.Interface()))
		return ffs
	}
}

func visitRoot(callback valueVisitor, root any, c *byteConsumer) {
	rootVal := reflect.ValueOf(root)

	path := valuePath{}
	if isPointerToSlice(rootVal) {
		visitRootSlice(callback, rootVal, c, path)
	} else {
		visitBreadthFirst(callback, rootVal, c, path)
	}

	//println("")
}

func isPointerToSlice(value reflect.Value) bool {
	return value.Kind() == reflect.Pointer && value.Elem().Kind() == reflect.Slice
}

func visitRootSlice(callback valueVisitor, pointerVal reflect.Value, c *byteConsumer, path valuePath) {
	path = path.add(pointerVal, "*")

	sliceVal := pointerVal.Elem()
	sliceType := sliceVal.Type().Elem()

	// Fill up the slice with all the available data
	for i := 0; c.len() > 0; i++ {
		// Create a new element for the slice
		pathName := fmt.Sprintf("[%d]", i)
		newVal := reflect.New(sliceType).Elem()

		// Fill in that new element with data
		visitBreadthFirst(callback, newVal, c, path.add(sliceVal, pathName))

		// Append the new element to the slice
		sliceVal.Set(reflect.Append(sliceVal, newVal))

		if !callback.canGrowRootSlice() {
			// If we don't make this check then the describer will
			// be unable to stop this slice from growing
			// indefinitely
			break
		}
	}
}

func visitBreadthFirst(callback valueVisitor, value reflect.Value, c *byteConsumer, path valuePath) {
	values := newDequeue[visitFunc]()

	visitFuncs := visitValue(callback, value, c, newEmptyFuzzTags(), path)
	values.addMany(visitFuncs)

	for values.len() != 0 {
		ff := values.popFirst()
		visitFuncs := ff()
		values.addMany(visitFuncs)
	}
}

func visitValue(callback valueVisitor, value reflect.Value, c *byteConsumer, tags fuzzTags, path valuePath) []visitFunc {
	if c.len() == 0 {
		// There are no more bytes to use to visit data
		return []visitFunc{}
	}

	switch value.Kind() {
	case reflect.Bool:
		callback.visitBool(value, c, tags, path)
		return []visitFunc{}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		callback.visitInt(value, c, tags, path)
		return []visitFunc{}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		callback.visitUint(value, c, tags, path)
		return []visitFunc{}

	case reflect.Uintptr:
		callback.visitUintptr(value, c, tags, path)
		return []visitFunc{}

	case reflect.Float32, reflect.Float64:
		callback.visitFloat(value, c, tags, path)
		return []visitFunc{}

	case reflect.Complex64, reflect.Complex128:
		callback.visitComplex(value, tags, path)
		return []visitFunc{}

	case reflect.Array:
		callback.visitArray(value, tags, path)

		newValues := []visitFunc{}
		for i := 0; i < value.Len(); i++ {
			pathVal := fmt.Sprintf("[%d]", i)
			newValues = append(newValues, visitValue(callback, value.Index(i), c, tags, path.add(value, pathVal))...)
		}
		return newValues

	case reflect.Chan:
		callback.visitChan(value, tags, path)
		return []visitFunc{}

	case reflect.Func:
		callback.visitFunc(value, tags, path)
		return []visitFunc{}

	case reflect.Interface:
		callback.visitInterface(value, tags, path)
		return []visitFunc{}

	case reflect.Map:
		mapLen := callback.visitMap(value, c, tags, path)

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
			newValues = append(newValues, visitValue(callback, mapKey, c, tags, path.add(value, "[key]"))...)

			// Create the value
			mapValP := reflect.New(valType)
			mapVal := mapValP.Elem()
			// Note here that the tags used to create this map are also
			// used to create the value
			newValues = append(newValues, visitValue(callback, mapVal, c, tags, path.add(value, "[value]"))...)

			// Add key/val to map
			//println("setting map element")
			value.SetMapIndex(mapKey, mapVal)
		}

		return newValues

	case reflect.Pointer:
		callback.visitPointer(value, c, tags, path)
		return []visitFunc{
			newVisitFunc(callback, value.Elem(), c, newEmptyFuzzTags(), path.add(value, "*")),
		}

	case reflect.Slice:
		sliceLen := callback.visitSlice(value, c, tags, path)

		newValues := []visitFunc{}
		for i := range sliceLen {
			pathVal := fmt.Sprintf("[%d]", i)
			newValues = append(newValues, visitValue(callback, value.Index(i), c, tags, path.add(value, pathVal))...)
		}
		return newValues

	case reflect.String:
		callback.visitString(value, c, tags, path)
		return []visitFunc{}

	case reflect.Struct:
		if !callback.visitStruct(value, tags, path) {
			// We allow the visitor to elect not to process a struct.
			// This was introduced to allow the Describe() function
			// to avoid infinite recursion
			return []visitFunc{}
		}

		if !value.CanSet() {
			// Can't set struct - ignore the struct and ignore its fields
			return []visitFunc{}
		}

		newValues := []visitFunc{}
		vType := value.Type()
		path = path.add(value, "("+vType.Name()+")")
		for i := 0; i < vType.NumField(); i++ {
			vField := value.Field(i)
			tField := vType.Field(i)
			tags := newFuzzTags(value, tField)
			if vField.CanSet() {
				newValues = append(newValues, visitValue(callback, vField, c, tags, path.add(value, tField.Name))...)
			} else {
				// We visit this unsettable field so we can describe it
				// It should not be filled and any returned fill functions are ignored
				visitValue(callback, vField, c, tags, path.add(value, tField.Name))
			}
		}

		return newValues

	case reflect.UnsafePointer:
		callback.visitUnsafePointer(value, tags, path)
		return []visitFunc{}

	default:
		panic(fmt.Errorf("unsupported kind %s", value.Kind()))
	}

	panic("unreachable")
}
