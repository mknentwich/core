package pdf

import "testing"

func TestFPDF(t *testing.T) {
	err := createBillFPDF(2345362)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
	}
}
