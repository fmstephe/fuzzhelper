package stack

import "fmt"

type Stack struct {
	top  int8 // This will cause bugs on overflow
	data []any
}

func New() *Stack {
	return &Stack{
		data: []any{},
	}
}

func (s *Stack) Push(val any) {
	s.top++
	if len(s.data) < int(s.top) {
		s.data = append(s.data, val)
	} else {
		s.data[s.top-1] = val
	}
}

func (s *Stack) Pop() any {
	if s.top == 0 {
		return nil
	}

	val := s.data[s.top-1]
	s.top--
	return val
}

func (s *Stack) String() string {
	return fmt.Sprintf("top: %d, data-length: %d", s.top, len(s.data))
}
