package alfred

import (
	"testing"
	"os"
	"io/ioutil"
)

var (
	testText = []byte("This is a test")
	testOutputPath = "test.txt"
)

func TestPrintResult(t *testing.T) {
	testString := string(testText)
	testOutput, err := os.Create(testOutputPath)
	if err != nil {
		t.Errorf("Create file - expected success, got error (%v).", err)
	}
	printResult(testString, testOutput)
	actualBytes, err := ioutil.ReadFile(testOutputPath)
	if err != nil {
		t.Errorf("Read output - expected success, got error (%v).", err)
	}
	actualOutput := string(actualBytes)
	if actualOutput != testString {
		t.Errorf("Expected output '%s', got '%s'.", testString, actualOutput)
	}
	if err := os.Remove(testOutputPath); err != nil {
		t.Errorf("Remove output file - expected success, got error (%v).", err)
	}
}
