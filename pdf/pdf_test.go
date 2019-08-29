package pdf

import "testing"

//Tests the generation of the bill pdf
func TestFPDF(t *testing.T) {
	err := createBillFPDF(2019082901)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
	}
}
