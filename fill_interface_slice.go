package fuzzhelper

type FillSliceStruct[T any] struct {
	allowableTypes []T `fuzz-ignore`
	Result         []T `fuzz-interface-method:"InterfaceOptions"`
}

func NewFillSliceStruct[T any](allowableTypes []T) *FillSliceStruct[T] {
	return &FillSliceStruct[T]{
		allowableTypes: allowableTypes,
	}
}

func (s *FillSliceStruct[T]) InterfaceOptions() []T {
	return s.allowableTypes
}

func FillSlice[T any](allowableTypes []T, bytes []byte) []T {
	fss := NewFillSliceStruct[T](allowableTypes)
	Fill(fss, bytes)
	return fss.Result
}
