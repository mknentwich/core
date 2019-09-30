package pdf

import (
	"os"
	"testing"
)

//Tests the generation of the bill pdf
func TestFPDF(t *testing.T) {
	insertTestData()
	f, err := os.OpenFile("example-bill.pdf", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
	}
	writee, err := writeBill(1)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
		return
	}
	err = writee(f)
	if err != nil {
		t.Errorf("Error on creating the bill pdf: %s", err.Error())
	}
}
