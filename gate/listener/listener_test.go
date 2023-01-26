package listener

import (
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveByIndex(t *testing.T) {
	chan1 := make(chan localmodel.Event)
	chan2 := make(chan localmodel.Event)
	chan3 := make(chan localmodel.Event)
	//chan4 := make(chan localmodel.Event)

	testCases := []struct {
		name              string
		originalChannels  []chan localmodel.Event
		resultingChannels []chan localmodel.Event
		indexToRemove     uint
		errorExpected     bool
	}{
		{
			name:              "single element",
			originalChannels:  []chan localmodel.Event{chan1},
			resultingChannels: []chan localmodel.Event{},
			indexToRemove:     0,
			errorExpected:     false,
		},
		{
			name:              "two elements #1",
			originalChannels:  []chan localmodel.Event{chan1, chan2},
			resultingChannels: []chan localmodel.Event{chan1},
			indexToRemove:     1,
			errorExpected:     false,
		},
		{
			name:              "two elements #1",
			originalChannels:  []chan localmodel.Event{chan1, chan2},
			resultingChannels: []chan localmodel.Event{chan2},
			indexToRemove:     0,
			errorExpected:     false,
		},
		{
			name:              "two elements error",
			originalChannels:  []chan localmodel.Event{chan1, chan2},
			resultingChannels: []chan localmodel.Event{chan2},
			indexToRemove:     2,
			errorExpected:     true,
		},
		{
			name:              "three elements #1",
			originalChannels:  []chan localmodel.Event{chan1, chan2, chan3},
			resultingChannels: []chan localmodel.Event{chan1, chan2},
			indexToRemove:     2,
			errorExpected:     false,
		},
		{
			name:              "three elements #2",
			originalChannels:  []chan localmodel.Event{chan1, chan2, chan3},
			resultingChannels: []chan localmodel.Event{chan1, chan3},
			indexToRemove:     1,
			errorExpected:     false,
		},
		{
			name:              "three elements #3",
			originalChannels:  []chan localmodel.Event{chan1, chan2, chan3},
			resultingChannels: []chan localmodel.Event{chan2, chan3},
			indexToRemove:     0,
			errorExpected:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := removeFromSliceByIndex(tc.originalChannels, tc.indexToRemove)
			if tc.errorExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.resultingChannels, result)
			}
		})
	}
}
