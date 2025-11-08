package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Set of tests for tag validation only. It's biggish so it gets its own file.

// Method tag slices are not allowed to be empty
type badEmptySlice struct {
	IntField int `fuzz-int-method:"EmptyOptions"`
}

func (b badEmptySlice) EmptyOptions() []int {
	return []int{}
}

// Method tag slices must have the correct type for the field
type badUint64ToUintAssignment struct {
	UintField uint `fuzz-int-method:"IntOptions"`
}

func (b badUint64ToUintAssignment) IntOptions() []uint64 {
	return []uint64{1}
}

// Method tag slices must have the correct type for the field
type badStringToIntAssignment struct {
	IntField int `fuzz-int-method:"StringOptions"`
}

func (b badStringToIntAssignment) StringOptions() []string {
	return []string{"this can't be assigned to int"}
}

// Method tag slices are not allowed to be nil
type badNilSliceOptions struct {
	IntField int `fuzz-int-method:"IntOptions"`
}

func (b badNilSliceOptions) IntOptions() []int {
	return nil
}

// Method tag methods must have the correct name
type badWrongMethodNameOptions struct {
	IntField int `fuzz-int-method:"IntOptions"`
}

func (b badWrongMethodNameOptions) WrongMethodName() []int {
	return []int{1, 2, 3}
}

type Satisfied interface {
	Satisfied() string
}

type valueStruct struct{}

func (v *valueStruct) Satisfied() string {
	return "Satisfies Satisfied interface"
}

// Method tag methods must return a slice of the exact interface type they will
// assign. Even if the underlying type can be assigned to the field - we only
// look at the interface type provided by the slice
type badWrongInterfaceForInterfaceValues struct {
	SatisfiedField Satisfied `fuzz-interface-method:"SatisfiedOptions"`
}

func (b badWrongInterfaceForInterfaceValues) SatisfiedOptions() []any {
	return []any{&valueStruct{}}
}

// Method tag values of interface types must _only_ be pointers, straight
// values are not allowed
type badValueTypeForInterfaceValues struct {
	AnyField any `fuzz-interface-method:"AnyOptions"`
}

func (b badValueTypeForInterfaceValues) AnyOptions() []any {
	return []any{valueStruct{}}
}

// Method tag methods must return a slice
type badMethodDoesNotReturnSlice struct {
	IntField int `fuzz-interface-method:"IntOptions"`
}

func (b badMethodDoesNotReturnSlice) IntOptions() int {
	return 1
}

// Method tag methods must return a slice
type badMethodRequiresArgument struct {
	IntField int `fuzz-interface-method:"IntOptions"`
}

func (b badMethodRequiresArgument) IntOptions(arg bool) []int {
	return []int{1, 2, 3}
}

// Method tag methods must return a slice
type badMethodNotExported struct {
	IntField int `fuzz-interface-method:"intOptions"`
}

func (b badMethodNotExported) intOptions() []int {
	return []int{1, 2, 3}
}

func TestFuzzTags_Bad(t *testing.T) {
	// The value of these bytes don't matter, but we do need _some_ bytes in order to reach the tags and expose the errors
	bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var testCases = []struct {
		name          string
		expectedError string
		value         any
	}{
		{
			name:          "empty slice",
			expectedError: "fuzzhelper.badEmptySlice.IntField has options method fuzzhelper.badEmptySlice.EmptyOptions(), but it returns an empty slice",
			value:         &badEmptySlice{},
		},
		{
			name:          "uint64 to uint assignment",
			expectedError: "fuzzhelper.badUint64ToUintAssignment.UintField cannot be assigned by every value returned by fuzzhelper.badUint64ToUintAssignment.IntOptions(), value of type uint64 cannot be assigned to uint",
			value:         &badUint64ToUintAssignment{},
		},
		{
			name:          "string to int assignment",
			expectedError: "fuzzhelper.badStringToIntAssignment.IntField cannot be assigned by every value returned by fuzzhelper.badStringToIntAssignment.StringOptions(), value of type string cannot be assigned to int",
			value:         &badStringToIntAssignment{},
		},
		{
			name:          "nil slice",
			expectedError: "fuzzhelper.badNilSliceOptions.IntField has options method fuzzhelper.badNilSliceOptions.IntOptions(), but it returns an empty slice",
			value:         &badNilSliceOptions{},
		},
		{
			name:          "wrong method name",
			expectedError: "fuzzhelper.badWrongMethodNameOptions.IntOptions() could not be found and can't be called",
			value:         &badWrongMethodNameOptions{},
		},
		{
			name:          "satisfies wrong interface for field",
			expectedError: "fuzzhelper.badWrongInterfaceForInterfaceValues.SatisfiedField cannot be assigned by every value returned by fuzzhelper.badWrongInterfaceForInterfaceValues.SatisfiedOptions(), value of type interface {} cannot be assigned to fuzzhelper.Satisfied",
			value:         &badWrongInterfaceForInterfaceValues{},
		},
		{
			name:          "value type for interface method (only pointers are allowed)",
			expectedError: "fuzzhelper.badValueTypeForInterfaceValues.AnyField cannot be assigned by every value returned by fuzzhelper.badValueTypeForInterfaceValues.AnyOptions(), value of type fuzzhelper.valueStruct must be a pointer (not a value type) to assign to interface type interface {}",
			value:         &badValueTypeForInterfaceValues{},
		},
		{
			name:          "method does not return slice",
			expectedError: "fuzzhelper.badMethodDoesNotReturnSlice.IntField cannot be assigned by every value returned by fuzzhelper.badMethodDoesNotReturnSlice.IntOptions(), method must return a slice, but it returns int",
			value:         &badMethodDoesNotReturnSlice{},
		},
		{
			name:          "method requires an argument",
			expectedError: "fuzzhelper.badMethodRequiresArgument.IntOptions() has arguments (found 1), must have no arguments",
			value:         &badMethodRequiresArgument{},
		},
		{
			name:          "method not exported",
			expectedError: "fuzzhelper.badMethodNotExported.intOptions() is not exported and can't be called",
			value:         &badMethodNotExported{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.PanicsWithError(t, testCase.expectedError, func() { Fill(testCase.value, bytes) })
		})
	}
}
