package fuzzhelper

import (
	"reflect"
)

var _ visitCallback = &fillVisitor{}

type fillVisitor struct {
}

func Fill(root any, c *ByteConsumer) {
	visitRoot(&fillVisitor{}, root, c)
}

func canSet(value reflect.Value) bool {
	// The initial value passed into Fill method must be an
	// instantiated struct/map/array/slice or a pointer to one of
	// these.  Once we drill past this unsettable level we will
	// fill in values recursively as we find them.
	if value.CanSet() {
		//println(": can set")
		return true
	}

	//println(": can't set")
	return false
}

func (v *fillVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		val := c.Uint64(BytesForNative)
		str := tags.stringValues[val%uint64(len(tags.stringValues))]

		//print("string ", len(str))
		if !canSet(value) {
			return []visitFunc{}
		}

		value.SetString(str)
		return []visitFunc{}
	}

	lengthVal := int(c.Int64(BytesForNative))
	strLength := tags.fitStringLength(lengthVal)

	//print("string ", strLength)
	if !canSet(value) {
		return []visitFunc{}
	}

	val := c.String(strLength)
	value.SetString(val)

	return []visitFunc{}
}

func (v *fillVisitor) visitBool(value reflect.Value, c *ByteConsumer, _ fuzzTags) []visitFunc {
	//print("bool")
	if !canSet(value) {
		return []visitFunc{}
	}
	val := c.Bool()
	value.SetBool(val)

	return []visitFunc{}
}

func (v *fillVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	//print("int")

	// First check there is a list of valid int values
	if len(tags.intValues) != 0 {
		val := c.Uint64(BytesForNative)
		intVal := tags.intValues[val%uint64(len(tags.intValues))]
		if !canSet(value) {
			return []visitFunc{}
		}

		value.SetInt(intVal)
		return []visitFunc{}
	}

	if !canSet(value) {
		return []visitFunc{}
	}
	val := c.Int64(value.Type().Size())
	fittedVal := tags.fitIntVal(val)
	value.SetInt(fittedVal)

	return []visitFunc{}
}

func (v *fillVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	//print("uint")

	// First check there is a list of valid uint values
	if len(tags.uintValues) != 0 {
		val := c.Uint64(BytesForNative)
		uintVal := tags.uintValues[val%uint64(len(tags.uintValues))]
		if !canSet(value) {
			return []visitFunc{}
		}

		value.SetUint(uintVal)
		return []visitFunc{}
	}

	if !canSet(value) {
		return []visitFunc{}
	}
	val := c.Uint64(value.Type().Size())
	fittedVal := tags.fitUintVal(val)
	value.SetUint(fittedVal)

	return []visitFunc{}
}

func (v *fillVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	//println("uintptr: ignored")
	return []visitFunc{}
}

func (v *fillVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	//print("float")

	// First check there is a list of valid uint values
	if len(tags.floatValues) != 0 {
		val := c.Uint64(BytesForNative)
		floatVal := tags.floatValues[val%uint64(len(tags.floatValues))]
		if !canSet(value) {
			return []visitFunc{}
		}

		value.SetFloat(floatVal)
		return []visitFunc{}
	}

	if !canSet(value) {
		return []visitFunc{}
	}

	val := c.Float64(value.Type().Size())
	fittedVal := tags.fitFloatVal(val)
	value.SetFloat(fittedVal)

	return []visitFunc{}
}

func (v *fillVisitor) visitStruct(value reflect.Value, c *ByteConsumer, _ fuzzTags) []visitFunc {
	//print("struct ", value.Type().Name())
	canSet(value)

	newValues := []visitFunc{}
	vType := value.Type()
	for i := 0; i < vType.NumField(); i++ {
		vField := value.Field(i)
		tField := vType.Field(i)
		tags := newFuzzTags(value, tField)
		newValues = append(newValues, visitValue(v, vField, c, tags)...)
	}

	return newValues
}

func (v *fillVisitor) visitPointer(value reflect.Value, c *ByteConsumer, _ fuzzTags) []visitFunc {
	//print("pointer")
	if !canSet(value) && value.IsNil() {
		return []visitFunc{}
	}

	if value.IsNil() {
		// If the value is nil - allocate a value for it to point to
		pType := value.Type()
		vType := pType.Elem()
		newVal := reflect.New(vType)
		value.Set(newVal)
	}
	return []visitFunc{
		newVisitFunc(v, value.Elem(), c, newEmptyFuzzTags()),
	}
}

func (v *fillVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags) int {
	val := int(c.Int64(BytesForNative))
	sliceLen := tags.fitSliceLengthVal(val)

	//print("slice ", sliceLen)
	if !canSet(value) && value.IsNil() {
		return 0
	}

	if value.IsNil() {
		newSlice := reflect.MakeSlice(value.Type(), sliceLen, sliceLen)
		value.Set(newSlice)
	}

	return sliceLen
}

// TODO there is a bug here where if the map cannot be set but is non-nil this function will try to set it
func (v *fillVisitor) visitMap(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	val := int(c.Int64(BytesForNative))
	mapLen := tags.fitMapLength(val)

	//print("map ", mapLen)
	if !canSet(value) && value.IsNil() {
		return []visitFunc{}
	}

	mapType := value.Type()
	keyType := mapType.Key()
	valType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)
	newValues := []visitFunc{}

	for range mapLen {
		// Create the key
		mapKeyP := reflect.New(keyType)
		mapKey := mapKeyP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the key
		newValues = append(newValues, visitValue(v, mapKey, c, tags)...)

		// Create the value
		mapValP := reflect.New(valType)
		mapVal := mapValP.Elem()
		// Note here that the tags used to create this map are also
		// used to create the value
		newValues = append(newValues, visitValue(v, mapVal, c, tags)...)

		// Add key/val to map
		//println("setting map element")
		newMap.SetMapIndex(mapKey, mapVal)
	}

	value.Set(newMap)

	return newValues
}

// TODO there is a bug here, if the channel can't be set, but is non-nil we will still try to set it
func (v *fillVisitor) visitChan(value reflect.Value, c *ByteConsumer, tags fuzzTags) []visitFunc {
	val := int(c.Int64(BytesForNative))
	chanLen := tags.fitChanLength(val)

	//print("chan ", chanLen)
	if !canSet(value) && value.IsNil() {
		return []visitFunc{}
	}

	chanType := value.Type()
	valType := chanType.Elem()

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), chanLen)
	newValues := []visitFunc{}

	for range chanLen {
		// Create an element for that channel
		newValP := reflect.New(valType)
		newVal := newValP.Elem()
		// Note here that the tags used to create this chan are also
		// used to create the values added to the channel
		newValues = append(newValues, visitValue(v, newVal, c, tags)...)

		// Put the element on the channel
		newValues = append(newValues, func() []visitFunc {
			newChan.Send(newVal)
			return []visitFunc{}
		})
	}

	value.Set(newChan)

	return newValues
}
