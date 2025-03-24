package fuzzhelper

import (
	"fmt"
	"os"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var _ valueVisitor = &describeVisitor{}

type describeVisitor struct {
}

func Describe(root any) {
	visitRoot(&describeVisitor{}, root, NewByteConsumer([]byte{1, 2, 3}))
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

func isExported(name string) bool {
	if name == "" {
		return false
	}

	firstRune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(firstRune)
}

func introDescription(value reflect.Value, tags fuzzTags, path valuePath) {
	fmt.Fprintf(os.Stdout, "%s\n", path.pathString(value))

	// Value is settable
	if value.CanSet() {
		return
	}

	// Field is not settable, lets find out why

	// If this is a field we want to make a very explicit message describing that fact
	if !isExported(tags.fieldName) {
		fmt.Fprintf(os.Stdout, "\tnot exported, will ignore\n")
		return
	}

	// If we reached here then the value cannot be set
	fmt.Fprintf(os.Stdout, "\tcan't set\n")
}

func (v *describeVisitor) visitBool(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)
	return
}

func (v *describeVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)

	if !value.CanSet() {
		// If we can't set this value don't provide any other details about it
		return
	}

	// First check if there is a list of valid string values
	if len(tags.intValues) != 0 {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("\tmethod (%s): %s", tags.intValuesMethod, methodValuesString(tags.intValues)))
		return
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.intRange.intMin, tags.intRange.intMax))
	return

}

func (v *describeVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)

	if !value.CanSet() {
		// If we can't set this value don't provide any other details about it
		return
	}

	// First check if there is a list of valid uint values
	if len(tags.uintValues) != 0 {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("\tmethod (%s): %s", tags.uintValuesMethod, methodValuesString(tags.uintValues)))
		return
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.uintRange.uintMin, tags.uintRange.uintMax))
	return
}

func (v *describeVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	notSupported(value, path)
}

func (v *describeVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)

	if !value.CanSet() {
		// If we can't set this value don't provide any other details about it
		return
	}

	// First check if there is a list of valid float values
	if len(tags.floatValues) != 0 {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("\tmethod (%s): %s", tags.floatValuesMethod, methodValuesString(tags.floatValues)))
		return
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %g max: %g", tags.floatRange.floatMin, tags.floatRange.floatMax))
	return
}

func (v *describeVisitor) visitComplex(value reflect.Value, tags fuzzTags, path valuePath) {
	// if this upsets you we can probably add it
	notSupported(value, path)
}

func (v *describeVisitor) visitArray(value reflect.Value, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)
}

func (v *describeVisitor) visitPointer(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	//introDescription(value, tags, path)

	if !value.CanSet() {
		return
	}

	// allocate a value for value to point to
	pType := value.Type()
	vType := pType.Elem()
	newVal := reflect.New(vType)
	value.Set(newVal)
}

func (v *describeVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) int {
	introDescription(value, tags, path)

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.sliceRange.uintRange.uintMin, tags.sliceRange.uintRange.uintMax))

	sliceLen := 1

	if !value.CanSet() {
		return 0
	}

	newSlice := reflect.MakeSlice(value.Type(), sliceLen, sliceLen)
	value.Set(newSlice)

	return sliceLen
}

// TODO there is a bug here where if the map cannot be set but is non-nil this function will try to set it
func (v *describeVisitor) visitMap(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) int {
	introDescription(value, tags, path)

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.mapRange.uintRange.uintMin, tags.mapRange.uintRange.uintMax))

	mapLen := 1

	if !value.CanSet() {
		return 0
	}

	mapType := value.Type()
	newMap := reflect.MakeMapWithSize(mapType, mapLen)
	value.Set(newMap)

	return mapLen
}

func (v *describeVisitor) visitChan(value reflect.Value, tags fuzzTags, path valuePath) {
	notSupported(value, path)
}

func (v *describeVisitor) visitFunc(value reflect.Value, tags fuzzTags, path valuePath) {
	notSupported(value, path)
}

func (v *describeVisitor) visitInterface(value reflect.Value, tags fuzzTags, path valuePath) {
	notSupported(value, path)
}

func (v *describeVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags, path valuePath) {
	introDescription(value, tags, path)

	if !value.CanSet() {
		// If we can't set this value don't provide any other details about it
		return
	}

	// First check if there is a list of valid string values
	if len(tags.stringValues) != 0 {
		fmt.Fprintf(os.Stdout, "\tmethod (%s): %s\n", tags.stringValuesMethod, methodValuesString(tags.stringValues))
		return
	}

	fmt.Fprintf(os.Stdout, "\trange min: %d max: %d\n", tags.stringRange.uintRange.uintMin, tags.stringRange.uintRange.uintMax)
	return
}

func (v *describeVisitor) visitStruct(value reflect.Value, tags fuzzTags, path valuePath) bool {
	recursion := path.containsType(value.Type())

	if !value.CanSet() || recursion {
		// We only describe a struct if we can't set it, or if it is recursive
		// If it can be set then it will be described via its fields
		introDescription(value, tags, path)
	}

	if recursion {
		fmt.Fprintf(os.Stdout, "\tRecursion...\n")
	}

	// If we've already visited (and described) a struct
	// we don't want to visit it again - so we return false
	return !recursion
}

func (v *describeVisitor) visitUnsafePointer(value reflect.Value, tags fuzzTags, path valuePath) {
	notSupported(value, path)
}

func notSupported(value reflect.Value, path valuePath) {
	fmt.Fprintf(os.Stdout, "%s\n", path.pathString(value))
	fmt.Fprintln(os.Stdout, "\tnot supported, will ignore")
}
