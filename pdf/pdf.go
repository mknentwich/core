package pdf

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"time"
)

//Numeric Layout constantes for right spacing inside the pdf
const (
	sizeHeader            = 10
	sizeText              = 11
	sizeBankData          = 14
	sizeTitle             = 18
	sizeLogo              = 26
	marginSide            = 25
	marginTop             = 12.5
	marginCustomerAddress = 50.8
	marginTitle           = marginCustomerAddress + 40
	marginArticles        = marginTitle + 25
	marginBankData        = marginArticles + 88
	marginQrCode          = marginBankData + 20
	marginQrCodeSide      = marginSide + 135
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

//amount, where a free delivery is granted to an order
var freeDelivery = 99.0

//Buffer, where PDF stream will be loaded into
var result bytes.Buffer

//Contains the bill information from database call
var billData OrderResultPDF
var billNumber string

//Info of the Company/Person
const sentenceTransfer = "Ich bitte Sie, den Betrag binnen 14 Tagen an das folgende Konto zu überweisen:"
const sentenceNoTaxes = "Dieser Betrag enthält keine Umsatzsteuer aufgrund §6(2)27 Kleinunternehmerregelung."

//Variables for svg handling and document print
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
		reference:  billNumber,
	}
}

//Initializes the ownAddress struct with static values
func initOwnAddress() *ownAddress {
	return &ownAddress{
		name:    "Markus Nentwich",
		street:  "Vereinsgasse 25/14",
		city:    "A-1020 Wien",
		phone:   "Telefon: +43699 / 10329882",
		email:   "E-Mail: kontakt@markus-nentwich.at",
		website: "Webseite: markus-nentwich.at",
	}
}

//Generates pdf bill from given orderId
//Returns reader, billNumber and error
func GeneratePDF(id int) (io.Reader, string, error) {
	billData, err = QueryOrderFromIdForPDF(id)
	if err != nil {
		return nil, "", err
	}
	billNumber = fmt.Sprint(billData.BillingDate) + fmt.Sprintf("%02d", billData.ReferenceCount)
	var address = initOwnAddress()
	var bank = initBankData()
	err = createBillPdf(*address, *bank)
	if err != nil {
		return nil, "", err
	}
	return bufio.NewReader(&result), billNumber, nil
}

//Creates the bill of an order as a pdf
func createBillPdf(address ownAddress, bank bankData) error {
	//Metadata and adding Page
	pdf := gofpdf.New("P", "mm", "A4", "resource/")
	pdf.SetTitle("Rechnung_"+billNumber, true)
	pdf.SetAuthor("Markus Nentwich", true)
	pdf.AddFont("wiener_melange", "B", "WienerMelange-Bold-new.json")
	pdf.AddPage()
	pdf.SetMargins(marginSide, marginTop, marginSide)
	translator = pdf.UnicodeTranslatorFromDescriptor("")
	width, _ := pdf.GetPageSize()
	cellWidthMax = width - 2*marginSide
	//MN Eigenverlag Image
	pdf.SetFont("wiener_melange", "B", sizeLogo)
	var titleTranslator = pdf.UnicodeTranslatorFromDescriptor("cp1252")
	pdf.SetX(marginSide)
	pdf.Cell(cellWidthMax, sizeLogo, titleTranslator("Nentwich Verlag"))
	//Own Address and Contact Info
	pdf.SetFont("Helvetica", "", sizeText)
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
	loadAddressesCustomer(pdf)
	//Title and Date of today
	pdf.SetX(marginSide)
	pdf.SetY(marginTitle)
	pdf.SetFontSize(sizeTitle)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax, sizeTitle, "Rechnung Nr. "+translator(billNumber), "", 0, "B", false, 0, "")
	pdf.SetFontSize(sizeText)
	pdf.SetFontStyle("")
	pdf.SetX(marginSide)
	pdf.CellFormat(cellWidthMax, sizeText, "Wien, am "+time.Now().Format("02.01.2006"), "", 1, "RT", false, 0, "")
	//Paint articles
	bank.price = loadArticles(pdf)
	//Paint bank data
	loadBankData(bank, pdf)
	return pdf.Output(&result)
}

//Paint the Addresses of the Customer
func loadAddressesCustomer(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.SetFontSize(sizeText)
	pdf.SetX(marginSide)
	pdf.SetY(marginCustomerAddress)

	if billData.Street != "" {
		pdf.SetFontStyle("b")
		pdf.Cell(cellWidthMax, sizeText, translator("Rechnungsadresse:"))
		pdf.SetFontStyle("")
		pdf.Ln(sizeText / 2)
	}
	if len(billData.Company) > 0 {
		pdf.Cell(cellWidthMax, sizeText, translator(billData.Company))
	} else {
		if len(billData.Salutation) > 0 {
			pdf.Cell(cellWidthMax, sizeText, translator(billData.Salutation))
		}
	}
	pdf.Ln(sizeText / 2)
	pdf.Cell(cellWidthMax, sizeText, translator(billData.FirstName+" "+billData.LastName))
	pdf.Ln(sizeText / 2)
	if billData.Street != "" {
		//use billing Address first
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.Street+" "+billData.BillingAddress.StreetNumber))
		pdf.Ln(sizeText / 2)
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.PostCode+" "+billData.BillingAddress.City))
		pdf.Ln(sizeText / 2)
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.Name))
		//paint deliveryAddress
		pdf.SetY(marginCustomerAddress)
		pdf.SetX(4.5 * marginSide)
		pdf.SetFontStyle("b")
		pdf.Cell(cellWidthMax, sizeText, translator("Lieferadresse:"))
		pdf.SetFontStyle("")
		pdf.Ln(sizeText / 2)
		pdf.SetX(4.5 * marginSide)
		pdf.Cell(cellWidthMax/2, sizeText, translator(billData.FirstName+" "+billData.LastName))
		pdf.Ln(sizeText / 2)
		pdf.SetX(4.5 * marginSide)
		pdf.Cell(cellWidthMax/2, sizeText, translator(billData.Street+" "+billData.StreetNumber))
		pdf.Ln(sizeText / 2)
		pdf.SetX(4.5 * marginSide)
		pdf.Cell(cellWidthMax/2, sizeText, translator(billData.PostCode+" "+billData.City))
		pdf.Ln(sizeText / 2)
		pdf.SetX(4.5 * marginSide)
		pdf.Cell(cellWidthMax/2, sizeText, translator(billData.Name))
	} else {
		//use billing Address only
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.Street+" "+billData.BillingAddress.StreetNumber))
		pdf.Ln(sizeText / 2)
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.PostCode+" "+billData.BillingAddress.City))
		pdf.Ln(sizeText / 2)
		pdf.Cell(cellWidthMax, sizeText, translator(billData.BillingAddress.Name))
	}
	return pdf
}

//Paint articles to the pdf, returns total price
//Checks, if free delivery is granted to this order
func loadArticles(pdf *gofpdf.Fpdf) float64 {
	price := float64(billData.ScoreAmount) * billData.Price
	pdf.SetX(marginSide)
	pdf.SetY(marginArticles)
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "Menge", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "Beschreibung", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Einzelpreis", "B", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Gesamtpreis", "B", 1, "R", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(cellWidthMax/9, sizeText, translator(fmt.Sprint(billData.ScoreAmount)), "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, translator(billData.Title), "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", billData.Price)+" €"), "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", price)+" €"), "", 0, "R", false, 0, "")
	pdf.Ln(sizeText / 2)
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "B", 0, "", false, 0, "")
	if billData.Street != "" {
		pdf.CellFormat(cellWidthMax/9*5, sizeText, translator("Versand ("+billData.Name+")"), "B", 0, "", false, 0, "")
	} else {
		pdf.CellFormat(cellWidthMax/9*5, sizeText, translator("Versand ("+billData.BillingAddress.Name+")"), "B", 0, "", false, 0, "")
	}
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "", "B", 0, "R", false, 0, "")
	if price >= freeDelivery {
		pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator("0.00 €"), "B", 1, "R", false, 0, "")
	} else {
		pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", billData.DeliveryPrice)+" €"), "B", 1, "R", false, 0, "")
		price += billData.DeliveryPrice
	}
	pdf.SetFontStyle("b")
	pdf.CellFormat(cellWidthMax/9, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*5, sizeText, "", "", 0, "", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, "Summe", "", 0, "R", false, 0, "")
	pdf.CellFormat(cellWidthMax/9*1.5, sizeText, translator(fmt.Sprintf("%.2f", price)+" €"), "", 1, "R", false, 0, "")
	pdf.SetFontStyle("")
	//Load no taxes sentence
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceNoTaxes), "", 0, "C", false, 0, "")
	return price
}

//Paint data of the bank connection
func loadBankData(bank bankData, pdf *gofpdf.Fpdf) {
	pdf.SetX(marginSide)
	pdf.SetY(marginBankData)
	pdf.CellFormat(cellWidthMax, sizeText, translator(sentenceTransfer), "", 1, "", false, 0, "")
	pdf.Ln(sizeText / 1.5)
	pdf.SetFontSize(sizeBankData)
	pdf.Cell(cellWidthMax/3, sizeText, "IBAN:")
	pdf.Cell(cellWidthMax/3, sizeText, translator(bank.ibanPretty))
	pdf.Ln(sizeText / 1.5)
	pdf.Cell(cellWidthMax/3, sizeText, "BIC:")
	pdf.Cell(cellWidthMax/3, sizeText, translator(bank.bic))
	pdf.Ln(sizeText / 1.5)
	pdf.Cell(cellWidthMax/3, sizeText, "Geldinstitut:")
	pdf.Cell(cellWidthMax/3, sizeText, translator(bank.institute))
	pdf.Ln(sizeText / 1.5)
	pdf.SetFontStyle("b")
	pdf.Cell(cellWidthMax/3, sizeText, "Zahlungsreferenz:")
	pdf.Cell(cellWidthMax/3, sizeText, translator(bank.reference))
	pdf.Ln(sizeText * 1.5)
	paintQRCode(generateQrCode(bank), pdf)
	pdf.SetFontSize(sizeText)
	pdf.CellFormat(cellWidthMax, sizeText, translator("Im Feld Zahlungsreferenz unbedingt obige Nummer anführen,"), "", 0, "C", false, 0, "")
	pdf.Ln(sizeText / 2)
	pdf.CellFormat(cellWidthMax, sizeText, translator("ansonsten kann die Zahlung nicht zugeordnet werden!"), "", 0, "C", false, 0, "")
	pdf.SetFontStyle("")
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
