package alfred

import (
	"encoding/json"
	"fmt"
	"io"
)

func printEntities(entities interface{}, w io.Writer) {
	bytes, err := json.Marshal(entities)
	if err != nil {
		panic(err)
	}
	printResult(string(bytes), w)
}

func printResult(result string, w io.Writer) {
	_, err := fmt.Fprint(w, result)
	if err != nil {
		panic(fmt.Errorf("Unable to print result '%s'. %v", result, err))
	}
}
