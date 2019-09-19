package pdf

import "testing"

//Tests the generation of the bill pdf
func TestFPDF(t *testing.T) {
	insertTestData()
	err := InitializePdfGeneration(1)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
	}
}
