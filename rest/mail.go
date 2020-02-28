package rest

import (
	"bytes"
	"encoding/base64"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/pdf"
	"github.com/mknentwich/core/utils"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"text/template"
	"time"
)

const (
	customerTemplate = "order-customer.mail.tmpl"
	ownerTemplate    = "order-owner.mail.tmpl"
)

type MailData struct {
	Order      *database.Order
	Date       string
	To         string
	From       string
	Attachment Attachment
	Boundary   string
}

type Attachment struct {
	Filename string
	Data     string
}

var customerMailBody *template.Template
var ownerMailBody *template.Template

func initializeTemplates() error {
	var err error
	ownerMailBody, err = utils.CreateTemplate(ownerTemplate)
	if err != nil {
		return err
	}
	customerMailBody, err = utils.CreateTemplate(customerTemplate)
	return err
}

func notify(order *database.Order) error {
	err := notifyCustomer(order)
	if err != nil {
		return err
	}
	return notifyOwner(order)
}

func notifyCustomer(order *database.Order) error {
	customerAddress := &mail.Address{Name: order.FirstName + " " + order.LastName, Address: order.Email}
	c := context.Conf
	data := &MailData{
		Order: order,
		Date:  time.Now().Format(time.RFC1123),
		To:    customerAddress.String(),
		From:  c.Mail.Address.String()}
	return sendMail(customerMailBody, []*mail.Address{customerAddress}, data)
}

func notifyOwner(order *database.Order) error {
	mails := make([]string, len(context.Conf.OrderRetrievers))
	for i, v := range context.Conf.OrderRetrievers {
		mails[i] = v.String()
	}
	pdfreader, filename, err := pdf.GeneratePDF(int(order.ID))
	if err != nil {
		return err
	}
	base64data, err := ioutil.ReadAll(pdfreader)
	if err != nil {
		return err
	}
	attachment := &Attachment{
		Filename: "Rechnung_" + filename + ".pdf",
		Data:     fold(base64.StdEncoding.EncodeToString(base64data), 76-1),
	}
	data := &MailData{
		Order:      order,
		Date:       time.Now().Format(time.RFC1123),
		To:         strings.Join(mails, ","),
		From:       context.Conf.Mail.Address.String(),
		Attachment: *attachment,
		Boundary:   "3809e216f78f1a242b12e913a8d3c6b0",
	}
	return sendMail(ownerMailBody, context.Conf.OrderRetrievers, data)
}

func fold(s string, n int) string {
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(s) - 1
	for i, rn := range s {
		buffer.WriteRune(rn)
		if i%n == n1 && i != l1 {
			buffer.WriteRune('\r')
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}

func sendMail(mailBody *template.Template, retriever []*mail.Address, data *MailData) error {
	mailCfg := context.Conf.Mail

	host, _, err := net.SplitHostPort(mailCfg.SMTPHost)
	if err != nil {
		return err
	}
	var msg []byte
	buf := bytes.NewBuffer(msg)
	err = mailBody.Execute(buf, data)
	msg = buf.Bytes()
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth(mailCfg.Address.String(), mailCfg.Username, mailCfg.Password, host)
	var to []string
	for _, retr := range retriever {
		to = append(to, retr.Address)
	}
	//fmt.Println(string(buf.Bytes()))
	return smtp.SendMail(mailCfg.SMTPHost, auth, mailCfg.Address.Address, to, msg)
}
