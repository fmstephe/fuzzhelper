package stack

import (
	"fmt"
	"testing"

	"github.com/fmstephe/fuzzhelper"
)

const (
	pushOpName = "push"
	popOpName  = "pop"
)

type stackOp interface {
	Name() string
	Value() string
}

type pushOp struct {
	PushValue string
}

func (o *pushOp) Name() string {
	return pushOpName
}

func (o *pushOp) Value() string {
	return o.PushValue
}

type popOp struct {
}

func (o *popOp) Name() string {
	return popOpName
}

func (o *popOp) Value() string {
	return ""
}

// Using StackFuzzStep we stress test the Stack with a random series of push/pop operations
func FuzzStack(f *testing.F) {
	f.Fuzz(func(t *testing.T, bytes []byte) {
		stack := New()

		// Run Describe, when running fuzzer as a normal test this will
		// help understanding how steps are being filled
		fuzzhelper.Describe(&[]stackOp{})

		// Construct the steps using the data in bytes
		steps := fuzzhelper.MakeSliceOf[stackOp]([]stackOp{&pushOp{}, &popOp{}}, bytes)
		count := 0

		defer func() {
			if p := recover(); p != nil {
				// recover a panic and add some extra information for debugging
				panic(fmt.Errorf("Panic stack{%s} steps: %d, %s", stack, count, p))
			}
		}()

		for _, step := range steps {
			count++
			switch step.Name() {
			case pushOpName:
				stack.Push(step.Value())
			case popOpName:
				stack.Pop()
			default:
				panic(fmt.Errorf("unknown operation: %d: %q", count, step.Name()))
			}
		}
	})
}
