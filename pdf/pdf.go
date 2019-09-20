package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"strings"
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
	marginQrCode          = marginBankData + 10
	marginQrCodeSide      = marginSide + 110
	marginNoTaxes         = marginBankData + 40
)

//Struct for contact details from the company
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
	iban       string
	ibanPretty string
	bic        string
	institute  string
	reference  string
	price      float64
}

//Contains the bill information from database call
var billData OrderResultPDF

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

//Initializes the bankdata struct
func initBankData() *bankData {
	return &bankData{
		iban:       "AT403209200000258475",
		ibanPretty: "AT40 3209 2000 0025 8475",
		bic:        "RLNWATWWGAE",
		institute:  "RAIFFEISEN-REGIONALBANK",
		reference:  "[referenceCount]",
	}
}

//Initializes the ownAddress struct with static values
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

//Generates pdf bill from given orderId
func InitializePdfGeneration(orderId int) error {
	billData = QueryOrderFromIdForPDF(orderId)
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
	pdf.SetTitle("RechnungsNr_[referenceCount]", true)
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
	//Paint billing address to PDF
	loadBillingAddress(pdf)
	//Title and Date of today
	pdf.SetX(marginSide)
	pdf.SetY(marginTitle)
	pdf.SetFontSize(sizeTitle)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax, sizeTitle, "Rechnung Nr. "+translator("[referenceCount]"), "", 0, "B", false, 0, "")
	pdf.SetFontSize(sizeText)
	pdf.SetFontStyle("")
	pdf.SetX(marginSide)
	pdf.CellFormat(cellWidthMax, sizeText, "Wien, am "+time.Now().Format("02.01.2006"), "", 1, "RT", false, 0, "")
	//Paint articles
	bank.price = loadArticles(pdf)
	//Paint bank data
	loadBankData(bank, pdf)
	//Load no taxes sentence
	pdf.SetX(marginSide)
	pdf.SetY(marginNoTaxes)
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceNoTaxes), "", 0, "", false, 0, "")
	err := pdf.OutputFileAndClose("example_bill.pdf")
	return err
}

//Paint the billing loadBillingAddress. All Static, just relevant for Spacing the document
func loadBillingAddress(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	var zHd string
	pdf.SetFontSize(sizeText)
	pdf.SetX(marginSide)
	pdf.SetY(marginCustomerAddress)
	if len(billData.Company) > 0 {
		zHd = "z.Hd. "
		pdf.Cell(cellWidthMax, sizeText, translator(billData.Company))
	} else {
		if len(billData.Salutation) > 0 {
			pdf.Cell(cellWidthMax, sizeText, translator(billData.Salutation))
		}
	}
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, translator(zHd+billData.FirstName+" "+billData.LastName))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, translator(billData.Street+" "+billData.StreetNumber))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, translator(billData.PostCode+" "+billData.City))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, translator(billData.State))
	return pdf
}

//Paint articles to the pdf, returns total price
func loadArticles(pdf *gofpdf.Fpdf) float64 {
	pdf.SetX(marginSide)
	pdf.SetY(marginArticles)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "Menge", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "Beschreibung", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Einzelpreis", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Gesamtpreis", "B", 1, "R", false, 0, "")
	price := float64(billData.ScoreAmount) * billData.Price
	pdf.SetFontStyle("")
	pdf.CellFormat(cellWidthMax/9, sizeText, translator(fmt.Sprint(billData.ScoreAmount)), "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, translator(billData.Title), "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", billData.Price)+" €"), "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", price)+" €"), "", 0, "R", false, 0, "")
	pdf.Ln(sizeText / 2)
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, translator("Versand ("+billData.State+")"), "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "", "B", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", checkDeliveryCosts())+" €"), "B", 1, "R", false, 0, "")
	price += checkDeliveryCosts()
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Summe", "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", price)+" €"), "", 1, "R", false, 0, "")
	pdf.SetFontStyle("")
	return price
}

//TODO: Implement database table for delivery costs to create dynamic values for each country
func checkDeliveryCosts() float64 {
	var price float64
	switch strings.ToLower(billData.State) {
	case "österreich":
		price = 3
	case "deutschland":
		price = 5.6
	}
	return price
}

//Paint data of the bank connection
func loadBankData(bank bankData, pdf *gofpdf.Fpdf) {
	pdf.SetX(marginSide)
	pdf.SetY(marginBankData)
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceTransfer), "", 1, "", false, 0, "")
	pdf.Cell(cellWidthMax/4, sizeText, "IBAN:")
	pdf.Cell(cellWidthMax/4, sizeText, translator(bank.ibanPretty))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "BIC:")
	pdf.Cell(cellWidthMax/4, sizeText, translator(bank.bic))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "Geldinstitut:")
	pdf.Cell(cellWidthMax/4, sizeText, translator(bank.institute))
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax/4, sizeText, "Zahlungsreferenz:")
	pdf.SetFontStyle("b")
	pdf.Cell(cellWidthMax/4, sizeText, translator(bank.reference))
	pdf.SetFontStyle("")
	pdf.Ln(sizeText / 2)
	paintQRCode(generateQrCode(bank), pdf)
}

//Paint generated qr code to bill
func paintQRCode(bitmap [][]bool, pdf *gofpdf.Fpdf) {
	var qrScale = 0.5
	for rows := range bitmap {
		for collums := range bitmap[rows] {
			if bitmap[rows][collums] == true {
				pdf.SetDrawColor(0, 0, 0)
				pdf.SetFillColor(0, 0, 0)
			} else {
				pdf.SetDrawColor(255, 255, 255)
				pdf.SetFillColor(255, 255, 255)
			}
			pdf.Rect(marginQrCodeSide+(float64(rows)*qrScale), marginQrCode+(float64(collums)*qrScale), qrScale, qrScale, "DF")
		}
	}
}
