package pdf

import (
	"github.com/jung-kurt/gofpdf"
	"strconv"
	"time"
)

//Numeric Layout constantes for right spacing inside the pdf
const (
	sizeHeader            = 10
	sizeText              = 12
	sizeTitle             = 18
	marginSide            = 25
	marginTop             = 12.5
	marginCustomerAddress = 50.8
	marginTitle           = 120
)

//Name of the used logo [svg file]
const svgFileName = "nentwich.svg"

var (
	sig gofpdf.SVGBasicType
	err error
)

//AddressInfo of the Company/Person
var ownAddress = []string{"Markus Nentwich", "Vereinsgasse 25/14", "A-1020 Wien", "Telefon: +43699 / 10329882", "E-Mail: nentwich94@gmx.at", "Webseite: "}

//Creates the bill of an order as a pdf
func createBillFPDF(referenceCount int) error {
	refCount := strconv.Itoa(referenceCount)

	//Metadata and adding Page
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", sizeText)
	pdf.SetTitle("RechnungsNr_"+refCount, true)
	pdf.SetAuthor("Nentwich Eigenverlag", true)
	pdf.AddPage()
	pdf.SetMargins(marginSide, marginTop, marginSide)

	translator := pdf.UnicodeTranslatorFromDescriptor("")
	width, _ := pdf.GetPageSize()
	cellWidthMax := width - 2*marginSide

	//MN Eigenverlag Image
	sig, err = gofpdf.SVGBasicFileParse(svgFileName)
	if err == nil {
		scale := 250 / sig.Wd
		scaleY := 75 / sig.Ht
		if scale > scaleY {
			scale = scaleY
		}
		pdf.SetLineCapStyle("round")
		pdf.SetLineWidth(0.26458332)
		pdf.SetDrawColor(238, 73, 87)
		pdf.SetXY(marginSide-5, marginTop)
		pdf.SVGBasicWrite(&sig, scale)
	} else {
		pdf.SetError(err)
	}

	//Own Address and Contact Info
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
	pdf.Cell(width, sizeText, "Salutation / [Company]")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "[z.Hd.] FirstName + LastName")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "BillingAddress 1")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "BillingAddress 2")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "STATE [if not austria]")

	//Title and Date of today
	pdf.SetFontSize(sizeTitle)
	pdf.SetFontStyle("b")
	pdf.SetX(marginSide)
	pdf.SetY(marginTitle)
	pdf.CellFormat(cellWidthMax, sizeTitle, "Rechnung Nr. "+refCount, "", 0, "", false, 0, "")
	pdf.SetFontSize(sizeText)
	pdf.SetFontStyle("")
	pdf.SetX(marginSide)
	pdf.CellFormat(cellWidthMax, sizeText, "Wien, am "+time.Now().Format("02.01.2006"), "", 1, "RB", false, 0, "")

	err := pdf.OutputFileAndClose("hello.pdf")
	return err
}
