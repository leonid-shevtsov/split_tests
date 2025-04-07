package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBias(t *testing.T) {
	examples := []struct {
		name   string
		input  string
		output []float64
		err    string
	}{
		{
			"happy case",
			"0=1,1=2.5",
			[]float64{1, 2.5},
			"",
		},
		{
			"bad format",
			"bad",
			nil,
			"not a valid bias declaration: bad",
		},
		{
			"bad pair",
			"0=1, bad",
			nil,
			"not a valid bias declaration:  bad",
		},
		{
			"bad index",
			"bad=0.5",
			nil,
			"failed to parse bias index: strconv.Atoi: parsing \"bad\": invalid syntax",
		},
		{
			"bad time",
			"0=bad",
			nil,
			"failed to parse bias time: strconv.ParseFloat: parsing \"bad\": invalid syntax",
		},
		{
			"index out of range",
			"3=0.5",
			nil,
			"bias index is not within the split number: 3",
		},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			output, err := parseBias(example.input, 2)
			if example.err != "" {
				require.EqualError(t, err, example.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, example.output, output)
			}
		})
	}
}
