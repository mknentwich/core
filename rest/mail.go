package rest

import (
	"bytes"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/utils"
	"net"
	"net/mail"
	"net/smtp"
	"text/template"
	"time"
)

const (
	mailDateFormat   = "Mon, 2 Jan 2006 15:04:05 +0700 "
	customerTemplate = "order-customer.mail.tmpl"
	ownerTemplate    = "order-owner.mail.tmpl"
)

type MailData struct {
	Order *database.Order
	Date  string
	To    string
	From  string
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
	data := &MailData{
		Order: order,
		Date:  time.Now().Format(mailDateFormat),
		To:    customerAddress.String(),
		From:  context.Conf.Mail.Address.String()}
	return sendMail(customerMailBody, []*mail.Address{customerAddress}, data)
}

func notifyOwner(order *database.Order) error {
	//TODO add data
	return sendMail(ownerMailBody, context.Conf.OrderRetrievers, nil)
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
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth(mailCfg.Address.String(), mailCfg.Username, mailCfg.Password, host)
	var to []string
	for _, retr := range retriever {
		to = append(to, retr.Address)
	}
	return smtp.SendMail(mailCfg.SMTPHost, auth, mailCfg.Address.Address, to, msg)
}
