package stack

import (
	"fmt"
	"testing"

	"github.com/fmstephe/fuzzhelper"
)

const (
	pushOp = "push"
	popOp  = "pop"
)

// StackFuzzStep is used to construct a series of push or pop operations
// for fuzzing test runs
type StackFuzzStep struct {
	Operation string `fuzz-string-method:"AllOperations"`
	PushValue string `fuzz-string-range:"1,10"`
}

// List of operations which we allow to execute at each step
func (s StackFuzzStep) AllOperations() []string {
	return []string{
		pushOp,
		popOp,
	}
}

// Using StackFuzzStep we stress test the Stack with a random series of push/pop operations
func FuzzStack(f *testing.F) {
	f.Fuzz(func(t *testing.T, bytes []byte) {
		stack := New()
		steps := &[]StackFuzzStep{}
		c := fuzzhelper.NewByteConsumer(bytes)
		fuzzhelper.Fill(steps, c)
		count := 0

		defer func() {
			if p := recover(); p != nil {
				// recover a panic and add some extra information for debugging
				panic(fmt.Errorf("Panic stack{%s} steps: %d, %s", stack, count, p))
			}
		}()

		for _, step := range *steps {
			count++
			switch step.Operation {
			case pushOp:
				stack.Push(step.PushValue)
			case popOp:
				stack.Pop()
			default:
				panic(fmt.Errorf("unknown operation: %d: %q", count, step.Operation))
			}
		}
	})
}
