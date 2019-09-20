package pdf

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

//Struct for qr code handling
type qrDataFields struct {
	serviceTag     string
	version        string
	coding         string
	function       string
	bic            string
	receiver       string
	iban           string
	amountCurrency string
	purpose        string
	reference      string
	text           string
	display        string
}

//Returns initialized qr code struct
func initializeQrDataFields(bank bankData) *qrDataFields {
	return &qrDataFields{
		serviceTag:     "BCD",
		version:        "001",
		coding:         "1",
		function:       "SCT",
		bic:            bank.bic,
		receiver:       "Markus Nentwich",
		iban:           bank.iban,
		amountCurrency: "EUR" + fmt.Sprintf("%.2f", bank.price),
		purpose:        "",
		reference:      bank.reference,
		text:           "",
		display:        "Ihre Transaktion an Nentwich Verlag wird vorbereitet",
	}
}

//Generates new QR Code from initialized qr code struct
//RecoveryLevel = Medium , VersionNumber = 13 due to qr code payment standards
func generateQrCode(bank bankData) [][]bool {
	var qrData = initializeQrDataFields(bank)
	var code *qrcode.QRCode
	code, _ = qrcode.New(qrData.serviceTag+"\n"+qrData.version+"\n"+qrData.coding+"\n"+qrData.function+
		"\n"+qrData.bic+"\n"+qrData.receiver+"\n"+qrData.iban+"\n"+qrData.amountCurrency+"\n"+
		qrData.purpose+"\n"+qrData.reference+"\n"+qrData.text+"\n"+qrData.display, qrcode.Medium)
	code.VersionNumber = 13
	return code.Bitmap()
}
