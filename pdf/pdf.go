package pdf

import (
	"github.com/jung-kurt/gofpdf"
	"strconv"
)

//Numeric Layout constantes for right spacing inside the pdf
const (
	sizeHeader            = 10
	sizeText              = 12
	sizeTitle             = 18
	marginSide            = 25
	marginTop             = 12.5
	marginCustomerAddress = 50.8
	marginTitle           = 100
)

//Name of the used logo [svg file]
const svgFileName = "nentwich.svg"

var (
	sig gofpdf.SVGBasicType
	err error
)

//AddressInfo of the Company/Person
var ownAddress = []string{"Nentwich Eigenverlag", "Markus Nentwich", "Adresse", "PLZ + Stadt", "Telefon: ", "E-Mail:", "Webseite: http://www.nentwichhh.at"}

//Creates the bill of an order as a pdf
func createBillFPDF(referenceCount int) error {
	refCount := strconv.Itoa(referenceCount)

	//Metadata and adding Page
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", sizeText)
	pdf.SetTitle("RechnungsNr_"+refCount, true)
	pdf.SetAuthor("Nentwich Eigenverlag", true)
	pdf.AddPage()

	translator := pdf.UnicodeTranslatorFromDescriptor("")
	width, _ := pdf.GetPageSize()
	cellWidthMax := width - 1.5*marginSide

	//MN Eigenverlag Image
	pdf.SetY(marginTop / 4)
	sig, err = gofpdf.SVGBasicFileParse(svgFileName)
	if err == nil {
		scale := 400 / sig.Wd
		scaleY := 120 / sig.Ht
		if scale > scaleY {
			scale = scaleY
		}
		pdf.SetLineCapStyle("round")
		pdf.SetLineWidth(0.26458332)
		pdf.SetDrawColor(238, 73, 87)
		pdf.SetXY(0, pdf.GetY())
		pdf.SVGBasicWrite(&sig, scale)
	} else {
		pdf.SetError(err)
	}

	//Own Address
	pdf.SetY(marginTop)
	pdf.SetFontSize(sizeHeader)
	for _, row := range ownAddress {
		pdf.SetX(marginSide)
		pdf.CellFormat(cellWidthMax, sizeHeader, translator(row), "", 0, "R", false, 0, "")
		pdf.Ln(sizeHeader / 2)
	}

	//Address of the customer, only for spacing reasons
	pdf.SetFontSize(sizeText)
	pdf.SetX(marginSide)
	pdf.SetY(marginCustomerAddress)
	pdf.Cell(cellWidthMax, sizeText, "Name")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "Company optional")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "address1")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "address2")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "address3")

	//Title
	pdf.SetFontSize(sizeTitle)
	pdf.SetFontStyle("b")
	pdf.SetX(marginSide)
	pdf.SetY(marginTitle)
	pdf.CellFormat(cellWidthMax, sizeTitle, "Rechnung Nr. "+refCount, "", 1, "", false, 0, "")

	err := pdf.OutputFileAndClose("hello.pdf")
	return err
}
