package alfred

import "os"

func PrintEntities(entities interface{}) {
	printEntities(entities, os.Stdout)
}

func PrintResult(result string) {
	printResult(result, os.Stdout)
}
