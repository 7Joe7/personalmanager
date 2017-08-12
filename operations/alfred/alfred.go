package alfred

import (
	"bufio"
	"encoding/json"
	"fmt"

	"io"

	"github.com/7joe7/personalmanager/resources"
)

type alfred struct{}

func NewAlfred() *alfred {
	return &alfred{}
}

func (a *alfred) PrintEntities(entities interface{}, writer io.Writer) {
	bytes, err := json.Marshal(entities)
	if err != nil {
		panic(err)
	}
	a.PrintResult(string(bytes), writer)
}

func (a *alfred) PrintResult(result string, writer io.Writer) {
	w := bufio.NewWriter(writer)
	_, err := w.Write([]byte(result))
	if err != nil {
		panic(fmt.Errorf("Unable to write result '%s'. %v", result, err))
	}
	_, err = w.Write([]byte(resources.STOP_CHARACTER))
	if err != nil {
		panic(fmt.Errorf("unable to write stop character. %v", err))
	}
	err = w.Flush()
	if err != nil {
		panic(fmt.Errorf("Unable to print result '%s'. %v", result, err))
	}
}
