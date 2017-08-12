package alfred

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Ahoj  string `json:"ahoj"`
	Behoj string `json:"behoj"`
}

func TestNewAlfred(t *testing.T) {
	NewAlfred()
}

func TestPrintResult(t *testing.T) {
	vectors := []struct {
		input  string
		output string
		err    error
	}{
		{
			input:  "This is a test",
			output: "This is a test",
			err:    nil,
		},
	}

	for _, v := range vectors {
		out := bytes.NewBuffer(make([]byte, 0, len(v.input)))
		a := &alfred{}
		a.PrintResult(v.input, out)
		o := make([]byte, len(v.output), len(v.output))
		_, err := out.Read(o)
		assert.Equal(t, v.err == nil, err == nil, "wrong failure condition")
		if err != nil && v.err != nil {
			assert.Equal(t, v.err.Error(), err.Error(), "wrong error")
			continue
		}
		if err != nil {
			continue
		}
		assert.Equal(t, v.output, string(o), "wrong output")
	}
}

func TestPrintEntities(t *testing.T) {
	vectors := []struct {
		input  interface{}
		output string
		err    error
	}{
		{
			input: &testObj{
				Ahoj:  "ahoj",
				Behoj: "behoj",
			},
			output: "{\"ahoj\":\"ahoj\",\"behoj\":\"behoj\"}",
			err:    nil,
		},
	}

	for _, v := range vectors {
		out := bytes.NewBuffer(make([]byte, 0))
		a := &alfred{}
		a.PrintEntities(v.input, out)
		o := make([]byte, len(v.output), len(v.output))
		_, err := out.Read(o)
		assert.Equal(t, v.err == nil, err == nil, "wrong failure condition")
		if err != nil && v.err != nil {
			assert.Equal(t, v.err.Error(), err.Error(), "wrong error")
			continue
		}
		if err != nil {
			continue
		}
		assert.Equal(t, v.output, string(o), "wrong output")
	}
}
