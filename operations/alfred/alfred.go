package alfred

import (
	"encoding/json"
	"fmt"
	"io"
)

type alfred struct {
	output io.Writer
}

func NewAlfred(o io.Writer) *alfred {
	return &alfred{o}
}

func (a *alfred) PrintEntities(entities interface{}) {
	bytes, err := json.Marshal(entities)
	if err != nil {
		panic(err)
	}
	a.PrintResult(string(bytes))
}

func (a *alfred) PrintResult(result string) {
	_, err := fmt.Fprint(a.output, result)
	if err != nil {
		panic(fmt.Errorf("Unable to print result '%s'. %v", result, err))
	}
}
