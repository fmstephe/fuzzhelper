package fuzzhelper

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
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

var pointerRegex = regexp.MustCompile(`\.(\**)\(`)

func pathString(value reflect.Value, path []string) string {
	pStr := strings.Join(path, ".")
	pStr = strings.ReplaceAll(pStr, "*.", "*")
	pStr = strings.ReplaceAll(pStr, ".[", "[")
	pStr = strings.ReplaceAll(pStr, ".(", "(")
	// This complex regex replace pulls leading pointer '*' inside the (typeName) parenthesis
	// This all feels a bit convoluted - will think about it a bit more
	pStr = pointerRegex.ReplaceAllString(pStr, "($1")

	if len(pStr) > 0 {
		pStr += " "
	}

	pStr = pStr + "(" + typeString(value.Type()) + ")"

	return pStr
}

func typeString(typ reflect.Type) string {
	switch typ.Kind() {
	case reflect.Pointer:
		return "*" + typeString(typ.Elem())
	case reflect.Slice:
		return "[]" + typeString(typ.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", typeString(typ.Key()), typeString(typ.Elem()))
	default:
		return typ.Name()
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

func isExported(name string) bool {
	if name == "" {
		return false
	}

	firstRune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(firstRune)
}

func introDescription(value reflect.Value, tags fuzzTags, path []string) {
	fmt.Fprintf(os.Stdout, "%s\n", pathString(value, path))

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

func (v *describeVisitor) visitBool(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
	return
}

func (v *describeVisitor) visitInt(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.intMin, tags.intMax))
	return

}

func (v *describeVisitor) visitUint(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.uintMin, tags.uintMax))
	return
}

func (v *describeVisitor) visitUintptr(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
	// Ignored
	return
}

func (v *describeVisitor) visitFloat(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %g max: %g", tags.floatMin, tags.floatMax))
	return
}

func (v *describeVisitor) visitArray(value reflect.Value, tags fuzzTags, path []string) {
	introDescription(value, tags, path)
}

func (v *describeVisitor) visitPointer(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
	//introDescription(value, tags, path)

	if !canSet(value) {
		return
	}

	// allocate a value for value to point to
	pType := value.Type()
	vType := pType.Elem()
	newVal := reflect.New(vType)
	value.Set(newVal)
}

func (v *describeVisitor) visitSlice(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) int {
	introDescription(value, tags, path)

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.sliceLengthMin, tags.sliceLengthMax))

	sliceLen := 1

	//fmt.Fprint(os.Stdout, "slice ", sliceLen)
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.mapLengthMin, tags.mapLengthMax))

	mapLen := 1

	//fmt.Fprint(os.Stdout, "map ", mapLen)
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

	fmt.Fprintln(os.Stdout, fmt.Sprintf("\trange min: %d max: %d", tags.chanLengthMin, tags.chanLengthMax))

	chanLen := 1

	//fmt.Fprint(os.Stdout, "chan ", chanLen)
	if !canSet(value) && value.IsNil() {
		return chanLen
	}

	// Create a channel
	newChan := reflect.MakeChan(value.Type(), chanLen)
	value.Set(newChan)

	return chanLen
}

func (v *describeVisitor) visitString(value reflect.Value, c *ByteConsumer, tags fuzzTags, path []string) {
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

	fmt.Fprintf(os.Stdout, "\trange min: %d max: %d\n", tags.stringLengthMin, tags.stringLengthMax)
	return
}

func (v *describeVisitor) visitStruct(value reflect.Value, tags fuzzTags, path []string) {
	if !value.CanSet() {
		// We only describe a struct if we can't set it
		// If it can be set then it will be described via its fields
		introDescription(value, tags, path)
	}
}
