package utils

import "testing"

//Fails a test if an unexpected error occur which has nothing to do with the test case.
func Unexpected(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Encountered unexpected error: %s", err.Error())
	}
}
