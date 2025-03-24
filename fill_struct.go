package fuzzhelper

import (
	"reflect"
	"strings"
)

var _ valueVisitor = &fillVisitor{}

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
func (v *fillVisitor) visitBool(value reflect.Value, c *ByteConsumer, _ fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//print("bool")
	if !canSet(value) {
		return
	}
	val := c.Bool()
	value.SetBool(val)

	return
}

func (v *fillVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//print("int")

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

	return
}

func (v *fillVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//print("uint")

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

	return
}

func (v *fillVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//println("uintptr: ignored")
	return
}

func (v *fillVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//print("float")

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

	return
}

func (v *fillVisitor) visitComplex(value reflect.Value, tags fuzzTags, path []string) {
	// Do nothing - complex numbers are simply not supported
	// we still visit them so we can _describe_ that we don't support them
	// We can add it, I just don't use them so I didn't bother at first
}

func (v *fillVisitor) visitArray(value reflect.Value, tags fuzzTags, path []string) {
	// Do nothing - the array is fixed in size and we do nothing here
	// Each of it's elements will be visited and we will fill those
}

func (v *fillVisitor) visitPointer(value reflect.Value, c *ByteConsumer, _ fuzzTags, path []string) {
	//print(leftPad(len(path)))
	//print("pointer")
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
}

func (v *fillVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	//print(leftPad(len(path)))
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
func (v *fillVisitor) visitMap(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	//print(leftPad(len(path)))
	val := int(c.Int64(BytesForNative))
	mapLen := tags.fitMapLength(val)

	//print("map ", mapLen)
	if !canSet(value) && value.IsNil() {
		return 0
	}

	mapType := value.Type()
	newMap := reflect.MakeMapWithSize(mapType, mapLen)
	value.Set(newMap)

	return mapLen
}

func (v *fillVisitor) visitChan(value reflect.Value, tags fuzzTags, path []string) {
	// Do nothing - channels are simply not supported
	// we still visit them so we can _describe_ that we don't support them
}

func leftPad(pad int) string {
	return strings.Repeat(" ", pad)
}

func (v *fillVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//print(leftPad(len(path)))
	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		val := c.Uint64(BytesForNative)
		str := tags.stringValues[val%uint64(len(tags.stringValues))]

		//print("string ", len(str))
		if !canSet(value) {
			return
		}

		value.SetString(str)
		return
	}

	lengthVal := int(c.Int64(BytesForNative))
	strLength := tags.fitStringLength(lengthVal)

	//print("string ", strLength)
	if !canSet(value) {
		return
	}

	val := c.String(strLength)
	value.SetString(val)

	return
}

func (v *fillVisitor) visitStruct(value reflect.Value, tags fuzzTags, path []string) {
	// Do nothing - the struct is fixed in size and we do nothing here
	// Each of it's fields will be visited and we will fill those
}
