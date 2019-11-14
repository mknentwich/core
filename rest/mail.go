package rest

import (
	"bytes"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"net"
	"net/mail"
	"net/smtp"
	"text/template"
)

const mailDateFormat = "Mon, 2 Jan 2006 15:04:05 +0700 "

type MailContent struct {
	Order    *database.Order
	Date     string
	Identity string
}

func notifyCustomer(order *database.Order) error {

}

func sendMail(mailBody template.Template, retriever []*mail.Address, data interface{}) error {
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
