package fuzzhelper

func ExampleDescribe_StringRange() {
	type testStruct struct {
		StringField string `fuzz-string-range:"1,5"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).StringField (string)
	//	range min: 1 max: 5
}

type stringMethodStruct struct {
	StringField string `fuzz-string-method:"StringValues"`
}

func (s *stringMethodStruct) StringValues() []string {
	return []string{
		"first",
		"second",
		"third",
		"fourth",
	}
}

func ExampleDescribe_StringMethod() {
	Describe(&stringMethodStruct{})
	// Output:*(stringMethodStruct).StringField (string)
	//	method (StringValues): [first second thi...
}

func ExampleDescribe_UnexportedString() {
	type testStruct struct {
		unexportedStringField string `fuzz-string-range:"1,5"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).unexportedStringField (string)
	//	not exported, will ignore
}

func ExampleDescribe_IntRange() {
	type testStruct struct {
		IntField int `fuzz-int-range:"-10,50"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).IntField (int)
	//	range min: -10 max: 50
}

type intMethodStruct struct {
	IntField int `fuzz-int-method:"IntValues"`
}

func (s *intMethodStruct) IntValues() []int64 {
	return []int64{
		-1,
		-2,
		-3,
		-4,
	}
}

func ExampleDescribe_IntMethod() {
	Describe(&intMethodStruct{})
	// Output:*(intMethodStruct).IntField (int)
	//	method (IntValues): [-1 -2 -3 -4]
}

func ExampleDescribe_UnexportedInt() {
	type testStruct struct {
		unexportedIntField int `fuzz-int-range:"-10,50"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).unexportedIntField (int)
	//	not exported, will ignore
}

func ExampleDescribe_UintRange() {
	type testStruct struct {
		UintField uint `fuzz-uint-range:"2,7"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).UintField (uint)
	//	range min: 2 max: 7
}

type uintMethodStruct struct {
	UintField uint `fuzz-uint-method:"UintValues"`
}

func (s *uintMethodStruct) UintValues() []uint64 {
	return []uint64{
		1,
		2,
		3,
		4,
	}
}

func ExampleDescribe_UintMethod() {
	Describe(&uintMethodStruct{})
	// Output:*(uintMethodStruct).UintField (uint)
	//	method (UintValues): [1 2 3 4]
}

func ExampleDescribe_UnexportedUint() {
	type testStruct struct {
		unexportedUintField uint `fuzz-uint-range:"2,7"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).unexportedUintField (uint)
	//	not exported, will ignore
}

func ExampleDescribe_FloatRange() {
	type testStruct struct {
		FloatField float64 `fuzz-float-range:"0.1,0.5"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).FloatField (float64)
	//	range min: 0.1 max: 0.5
}

type float64MethodStruct struct {
	FloatField float64 `fuzz-float-method:"FloatValues"`
}

func (s *float64MethodStruct) FloatValues() []float64 {
	return []float64{
		0.01,
		0.02,
		0.03,
		0.04,
	}
}

func ExampleDescribe_FloatMethod() {
	Describe(&float64MethodStruct{})
	// Output:*(float64MethodStruct).FloatField (float64)
	//	method (FloatValues): [0.01 0.02 0.03 0...
}

func ExampleDescribe_UnexportedFloat() {
	type testStruct struct {
		unexportedFloatField float64 `fuzz-float-range:"0.1,0.5"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).unexportedFloatField (float64)
	//	not exported, will ignore
}

func ExampleDescribe_SliceRange() {
	type testStruct struct {
		SliceField []float64 `fuzz-slice-range:"3,20" fuzz-float-range:"0.2,0.7"`
	}

	Describe(&testStruct{})
	// Output:*(testStruct).SliceField ([]float64)
	//	range min: 3 max: 20
	//*(testStruct).SliceField[0] (float64)
	//	range min: 0.2 max: 0.7
}

type parentStruct struct {
	// struct processed fifth
	PointerPointerChild **childStruct
	// struct processed first
	unexportedChild childStruct
	// struct processed fourth
	PointerChild *childStruct
	// struct processed second
	ValueChild childStruct
	// slice processed third, but the elements in the slice (via a pointer)
	// are processed sixth
	SliceChild []*childStruct
}

type childStruct struct {
	BoolField   bool
	StringField string
}

func ExampleDescribe_StructInStruct() {
	// Take careful note that because we defer processing of pointer
	// values, i.e. PointerChild and PointerPointerChild The value fields,
	// unexportedChild and ValueChild will be processed first, and
	// therefore described first.
	Describe(&parentStruct{})
	// Output:*(parentStruct).unexportedChild (childStruct)
	//	not exported, will ignore
	//*(parentStruct).ValueChild(childStruct).BoolField (bool)
	//*(parentStruct).ValueChild(childStruct).StringField (string)
	//	range min: 0 max: 20
	//*(parentStruct).SliceChild ([]*childStruct)
	//	range min: 0 max: 20
	//*(parentStruct).PointerChild(*childStruct).BoolField (bool)
	//*(parentStruct).PointerChild(*childStruct).StringField (string)
	//	range min: 0 max: 20
	//*(parentStruct).SliceChild[0](*childStruct).BoolField (bool)
	//*(parentStruct).SliceChild[0](*childStruct).StringField (string)
	//	range min: 0 max: 20
	//*(parentStruct).PointerPointerChild(**childStruct).BoolField (bool)
	//*(parentStruct).PointerPointerChild(**childStruct).StringField (string)
	//	range min: 0 max: 20
}
