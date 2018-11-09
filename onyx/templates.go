package onyx

import (
	"io/ioutil"
	"testing"
)

// Compares a string
// with the contents of the
// specified file.
func EqualToFile(str string, filePath string) bool {

	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	if str != string(fileBytes) {
		return false
	}

	return true

}

// Detect if rendering template caused
// a panic. This function should be invoked
// with keyword `defer`
func DetectPanic(t *testing.T, templateName string) {
	if n := recover(); n != nil {
		t.Errorf("Test for %s failed!", templateName)
	}
}
