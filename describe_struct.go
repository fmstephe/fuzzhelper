package fuzzhelper

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

var _ valueVisitor = &describeVisitor{}

type describeVisitor struct {
}

func Describe(root any) {
	visitRoot(&describeVisitor{}, root, NewByteConsumer([]byte{1, 2, 3}))
}

func pathString(value reflect.Value, path []string) string {
	pStr := ""

	for i, p := range path {
		if strings.HasPrefix(p, "[") && len(pStr) > 0 {
			// Delete previous "."
			pStr = pStr[:len(pStr)-1]
		}

		pStr += p
		if !(p == "*" || i == len(path)-1) {
			pStr += "."
		}
	}

	if len(pStr) > 0 {
		pStr += " "
	}

	pStr = pStr + "(" + kindString(value) + ")"

	return pStr
}

func kindString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Pointer:
		return "*"
	case reflect.Struct:
		return value.Type().Name()
	default:
		return value.Kind().String()
	}
}

func shortenString(s string) string {
	const limit = 20
	if len(s) <= limit {
		return s
	}
	return s[:limit-3] + "..."
}

func methodValuesString[T any](v []T) string {
	vStr := ""
	if len(v) == 0 {
		return vStr
	}
	return shortenString(fmt.Sprintf("%v", v))
}

func introDescription(value reflect.Value, tags fuzzTags, path []string) {
	println(pathString(value, path))
	print(leftPad(len(path)))
	if value.Kind() == reflect.Pointer {
		print("*")
	} else {
		print(value.Kind().String())
	}

	// Value is settable
	if value.CanSet() {
		println(": can set")
		return
	}

	// If this is a field we want to make a very explicit message describing that fact
	if tags.fieldName != "" {
		firstRune, _ := utf8.DecodeRuneInString(tags.fieldName)
		if !unicode.IsUpper(firstRune) {
			println(": not exported, will ignore")
			return
		}
	}

	// If we reached here then the value cannot be set
	println(": can't set")
}

func (v *describeVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)

	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		print(leftPad(len(path)))
		println(fmt.Sprintf("method (%s): %s", tags.stringValuesMethod, methodValuesString(tags.stringValues)))
		return
	}

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.stringLengthMin, tags.stringLengthMax))
	return

}

func (v *describeVisitor) visitBool(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
	return
}

func (v *describeVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)

	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		print(leftPad(len(path)))
		println(fmt.Sprintf("method (%s): %s", tags.intValuesMethod, methodValuesString(tags.intValues)))
		return
	}

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.intMin, tags.intMax))
	return

}

func (v *describeVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
	//print("uint")

	// First check if there is a list of valid uint values
	if len(tags.uintValues) != 0 {
		print(leftPad(len(path)))
		println(fmt.Sprintf("method (%s): %s", tags.uintValuesMethod, methodValuesString(tags.uintValues)))
		return
	}

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.uintMin, tags.uintMax))
	return
}

func (v *describeVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
	// Ignored
	return
}

func (v *describeVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)

	// First check if there is a list of valid float values
	if len(tags.floatValues) != 0 {
		print(leftPad(len(path)))
		println(fmt.Sprintf("method (%s): %s", tags.floatValuesMethod, methodValuesString(tags.floatValues)))
		return
	}

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %f max: %f", tags.floatMin, tags.floatMax))
	return
}

func (v *describeVisitor) visitArray(value reflect.Value, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
}

func (v *describeVisitor) visitPointer(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
}

func (v *describeVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	introDescription(value, tags, path)

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.sliceLengthMin, tags.sliceLengthMax))

	sliceLen := 1

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
func (v *describeVisitor) visitMap(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	introDescription(value, tags, path)

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.mapLengthMin, tags.mapLengthMax))

	mapLen := 1

	//print("map ", mapLen)
	if !canSet(value) && value.IsNil() {
		return 0
	}

	mapType := value.Type()
	newMap := reflect.MakeMapWithSize(mapType, mapLen)
	value.Set(newMap)

	return mapLen
}

// TODO there is a bug here, if the channel can't be set, but is non-nil we will still try to set it
func (v *describeVisitor) visitChan(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	introDescription(value, tags, path)

	print(leftPad(len(path)))
	println(fmt.Sprintf("range min: %d max: %d", tags.chanLengthMin, tags.chanLengthMax))

	chanLen := 1

	//print("chan ", chanLen)
	if !canSet(value) && value.IsNil() {
		return chanLen
	}

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), chanLen)
	value.Set(newChan)

	return chanLen
}
