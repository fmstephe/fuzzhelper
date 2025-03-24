package fuzzhelper

import (
	"reflect"
)

var _ valueVisitor = &fillVisitor{}

type fillVisitor struct {
}

func Fill(root any, c *ByteConsumer) {
	visitRoot(&fillVisitor{}, root, c)
}

func (v *fillVisitor) canGrowRootSlice() bool {
	return true
}

func (v *fillVisitor) visitBool(value reflect.Value, c *ByteConsumer, _ fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//print("bool")
	if !value.CanSet() {
		return
	}

	val := c.Bool()
	value.SetBool(val)
}

func (v *fillVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//print("int")
	if !value.CanSet() {
		return
	}

	// First check there is a list of valid int values
	if tags.intValues.wasSet {
		val := c.Uint64(BytesForNative)
		intVal := tags.intValues.value[val%uint64(len(tags.intValues.value))]

		value.SetInt(intVal)
		return
	}

	val := c.Int64(value.Type().Size())
	fittedVal := tags.intRange.fit(val)
	value.SetInt(fittedVal)
}

func (v *fillVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//print("uint")
	if !value.CanSet() {
		return
	}

	// First check there is a list of valid uint values
	if tags.uintValues.wasSet {
		val := c.Uint64(BytesForNative)
		uintVal := tags.uintValues.value[val%uint64(len(tags.uintValues.value))]

		value.SetUint(uintVal)
		return
	}

	val := c.Uint64(value.Type().Size())
	fittedVal := tags.uintRange.fit(val)
	value.SetUint(fittedVal)
}

func (v *fillVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//println("uintptr: ignored")
}

func (v *fillVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//print("float")
	if !value.CanSet() {
		return
	}

	// First check there is a list of valid uint values
	if tags.floatValues.wasSet {
		val := c.Uint64(BytesForNative)
		floatVal := tags.floatValues.value[val%uint64(len(tags.floatValues.value))]
		value.SetFloat(floatVal)
		return
	}

	val := c.Float64(value.Type().Size())
	fittedVal := tags.floatRange.fit(val)
	value.SetFloat(fittedVal)
}

func (v *fillVisitor) visitComplex(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - complex numbers are simply not supported
	// we still visit them so we can _describe_ that we don't support them
	// We can add it, I just don't use them so I didn't bother at first
}

func (v *fillVisitor) visitArray(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - the array is fixed in size and we do nothing here
	// Each of it's elements will be visited and we will fill those
}

func (v *fillVisitor) visitPointer(value reflect.Value, c *ByteConsumer, _ fuzzTags, path valuePath) {
	//print(leftPad(len(path)))
	//print("pointer")
	if !value.CanSet() {
		return
	}

	// If the value is nil - allocate a value for it to point to
	pType := value.Type()
	vType := pType.Elem()
	newVal := reflect.New(vType)
	value.Set(newVal)
}

func (v *fillVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) int {
	//print(leftPad(len(path)))
	val := int(c.Int64(BytesForNative))
	sliceLen := tags.sliceRange.fit(val)

	//print("slice ", sliceLen)
	if !value.CanSet() {
		return 0
	}

	newSlice := reflect.MakeSlice(value.Type(), sliceLen, sliceLen)
	value.Set(newSlice)

	return sliceLen
}

// TODO there is a bug here where if the map cannot be set but is non-nil this function will try to set it
func (v *fillVisitor) visitMap(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) int {
	//print(leftPad(len(path)))
	val := int(c.Int64(BytesForNative))
	mapLen := tags.mapRange.fit(val)

	//print("map ", mapLen)
	if !value.CanSet() {
		return 0
	}

	mapType := value.Type()
	newMap := reflect.MakeMapWithSize(mapType, mapLen)
	value.Set(newMap)

	return mapLen
}

func (v *fillVisitor) visitChan(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - channels are simply not supported
	// we still visit them so we can _describe_ that we don't support them
}

func (v *fillVisitor) visitFunc(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - functions are simply not supported
	// we still visit them so we can _describe_ that we don't support them
}

func (v *fillVisitor) visitInterface(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - interfaces are simply not supported
	// we still visit them so we can _describe_ that we don't support them
}

func (v *fillVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	if !value.CanSet() {
		return
	}

	// First check if there is a list of valid string values
	if tags.stringValues.wasSet {
		val := c.Uint64(BytesForNative)
		str := tags.stringValues.value[val%uint64(len(tags.stringValues.value))]

		value.SetString(str)
		return
	}

	lengthVal := int(c.Int64(BytesForNative))
	strLength := tags.stringRange.fit(lengthVal)

	val := c.String(strLength)
	value.SetString(val)
}

func (v *fillVisitor) visitStruct(value reflect.Value, tags fuzzTags, path valuePath) bool {
	// Do nothing - the struct is fixed in size and we do nothing here
	// Each of it's fields will be visited and we will fill those
	return true
}

func (v *fillVisitor) visitUnsafePointer(value reflect.Value, tags fuzzTags, path valuePath) {
	// Do nothing - unsafe pointers are simply not supported
	// we still visit them so we can _describe_ that we don't support them
}
