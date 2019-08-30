package pdf

import (
	"github.com/jung-kurt/gofpdf"
	"time"
)

//Numeric Layout constantes for right spacing inside the pdf
const (
	sizeHeader            = 10
	sizeText              = 11
	sizeTitle             = 18
	marginSide            = 25
	marginTop             = 12.5
	marginCustomerAddress = 50.8
	marginTitle           = marginCustomerAddress + 40
	marginArticles        = marginTitle + 25
	marginBankData        = marginArticles + 80
	marginNoTaxes         = marginBankData + 40
)

type ownAddress struct {
	name    string
	street  string
	city    string
	phone   string
	email   string
	website string
}

//Struct for the bankdata
type bankData struct {
	iban      string
	bic       string
	institute string
	reference string
}

//Info of the Company/Person
const sentenceTransfer = "Ich bitte Sie, den Betrag binnen 14 Tagen an das folgende Konto zu überweisen:"
const sentenceNoTaxes = "Dieser Betrag enthält keine Umsatzsteuer aufgrund §6(2)27 Kleinunternehmerregelung."

//Name of the used logo [svg file]
const svgEigenverlag = "nentwich.svg"

//Variables for svg handling and document print
var sig gofpdf.SVGBasicType
var err error
var translator func(string) string
var cellWidthMax float64

func initBankData() *bankData {
	return &bankData{
		iban:      "AT40 3209 2000 0025 8475",
		bic:       "RLNWATWWGAE",
		institute: "RAIFFEISEN-REGIONALBANK",
		reference: "",
	}
}

func initOwnAddress() *ownAddress {
	return &ownAddress{
		name:    "Markus Nentwich",
		street:  "Vereinsgasse 25/14",
		city:    "A-1020 Wien",
		phone:   "Telefon: +43699 / 10329882",
		email:   "E-Mail: nentwich94@gmx.at",
		website: "Webseite: ",
	}
}

func initializePdfGeneration() error {
	var address = initOwnAddress()
	var bank = initBankData()
	err := createBillPdf(*address, *bank)
	return err
}

//Creates the bill of an order as a pdf
func createBillPdf(address ownAddress, bank bankData) error {
	//Metadata and adding Page
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", sizeText)
	pdf.SetTitle("RechnungsNr_referenceCount", true)
	pdf.SetAuthor("Nentwich Eigenverlag", true)
	pdf.AddPage()
	pdf.SetMargins(marginSide, marginTop, marginSide)
	translator = pdf.UnicodeTranslatorFromDescriptor("")
	width, _ := pdf.GetPageSize()
	cellWidthMax = width - 2*marginSide
	//MN Eigenverlag Image
	sig, err = gofpdf.SVGBasicFileParse(svgEigenverlag)
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
	pdf.SetX(marginSide)
	pdf.SetY(marginTop)
	pdf.SetFontSize(sizeHeader)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.name), "", 0, "R", false, 0, "")
	pdf.Ln(sizeHeader / 2)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.street), "", 0, "R", false, 0, "")
	pdf.Ln(sizeHeader / 2)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.city), "", 0, "R", false, 0, "")
	pdf.Ln(sizeHeader / 2)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.phone), "", 0, "R", false, 0, "")
	pdf.Ln(sizeHeader / 2)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.email), "", 0, "R", false, 0, "")
	pdf.Ln(sizeHeader / 2)
	pdf.CellFormat(cellWidthMax, sizeHeader, translator(address.website), "", 0, "R", false, 0, "")

	//Static load
	loadBillingAddress(pdf)
	//Title and Date of today
	pdf.SetX(marginSide)
	pdf.SetY(marginTitle)
	pdf.SetFontSize(sizeTitle)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax, sizeTitle, "Rechnung Nr. "+translator("referenceCount"), "", 0, "B", false, 0, "")
	pdf.SetFontSize(sizeText)
	pdf.SetFontStyle("")
	pdf.SetX(marginSide)
	pdf.CellFormat(cellWidthMax, sizeText, "Wien, am "+time.Now().Format("02.01.2006"), "", 1, "RT", false, 0, "")
	//Static load
	loadArticles(pdf)
	//Load bank data
	loadBankData(bank, pdf)
	//Load no Taxes and Signature
	pdf.SetX(marginSide)
	pdf.SetY(marginNoTaxes)
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceNoTaxes), "", 0, "", false, 0, "")

	err := pdf.OutputFileAndClose("example_bill.pdf")
	return err
}

//Paint the billing loadBillingAddress. All Static, just relevant for Spacing the document
func loadBillingAddress(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.SetFontSize(sizeText)
	pdf.SetX(marginSide)
	pdf.SetY(marginCustomerAddress)
	pdf.Cell(cellWidthMax, sizeText, "Salutation / [Company]")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "[z.Hd.] FirstName + LastName")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "BillingAddress 1")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "BillingAddress 2")
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, "STATE [if not austria]")
	return pdf
}

//Paint Articles. All Static, just relevant for Spacing the document
func loadArticles(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	//Header
	pdf.SetX(marginSide)
	pdf.SetY(marginArticles)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "Menge", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "Beschreibung", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Einzelpreis", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Gesamtpreis", "B", 1, "R", false, 0, "")
	//Body
	pdf.SetFontStyle("")
	pdf.CellFormat(cellWidthMax/9, sizeText, "1", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "Eine letzte Runde (Blasorchesterfassung)", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator("39,00 €"), "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator("39,00 €"), "", 0, "R", false, 0, "")
	pdf.Ln(sizeText / 2)
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, translator("Versand (Österreich)"), "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "", "B", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator("3,00 €"), "B", 1, "R", false, 0, "")
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Summe", "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator("42,00 €"), "", 1, "R", false, 0, "")
	pdf.SetFontStyle("")
	return pdf
}

//Paint data of the bank connection
func loadBankData(bank bankData, pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.SetX(marginSide)
	pdf.SetY(marginBankData)
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceTransfer), "", 1, "", false, 0, "")
	pdf.Cell(cellWidthMax/4, sizeText, "IBAN:")
	pdf.Cell(cellWidthMax/4, sizeText, bank.iban)
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "BIC:")
	pdf.Cell(cellWidthMax/4, sizeText, bank.bic)
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "Geldinstitut:")
	pdf.Cell(cellWidthMax/4, sizeText, bank.institute)
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "Verwendungszweck:")
	pdf.Cell(cellWidthMax/4, sizeText, bank.reference)
	pdf.Ln(sizeText / 2)
	return pdf
}
