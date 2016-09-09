package alfred

import (
	"fmt"
	"encoding/json"
	"log"
)

func printEntities(entities interface{}) {
	bytes, err := json.Marshal(entities)
	if err != nil {
		log.Fatalf("Unable to marshal entities. %v", err)
	}
	printResult(string(bytes))
}

func printResult(result string) {
	_, err := fmt.Print(result)
	if err != nil {
		log.Fatalf("Unable to print result '%s'. %v", result, err)
	}
}