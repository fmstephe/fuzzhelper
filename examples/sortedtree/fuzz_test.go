package sortedtree

import (
	"slices"
	"testing"

	"github.com/fmstephe/fuzzhelper"
	"github.com/stretchr/testify/assert"
)

type SortedTreeFuzzStep struct {
	Value    string `fuzz-string-range:"1,256"`
	NextStep *SortedTreeFuzzStep
}

// Using StackFuzzStep we stress test the SortedTree with a random series of added string values
//
// We verify that the sorting is working correctly by
func FuzzSortedTree(f *testing.F) {
	f.Fuzz(func(t *testing.T, bytes []byte) {
		tree := New()
		step := &SortedTreeFuzzStep{}
		c := fuzzhelper.NewByteConsumer(bytes)
		fuzzhelper.Fill(step, c)

		sortedStrings := []string{}

		for step != nil {
			tree.Add(step.Value)
			sortedStrings = append(sortedStrings, step.Value)
			slices.Sort(sortedStrings)
			assert.Equal(t, sortedStrings[0], tree.Least(), "Tree %s", tree)
			assert.Equal(t, sortedStrings[len(sortedStrings)-1], tree.Greatest(), "Tree %s", tree)
			step = step.NextStep
		}
	})
}
