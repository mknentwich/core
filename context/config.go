package context

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/mail"
	"os"
)

const configFile = "config.json"

var Conf Configuration

//Struct for the configuration of the application.
type Configuration struct {
	Authentication       bool             `json:"authentication"`
	DavPaths             DavPaths         `json:"davPaths"`
	GeneratedDirectory   string           `json:"generatedDirectory"`
	Host                 string           `json:"host"`
	JWTExpirationMinutes int              `json:"jwtExpirationMinutes"`
	JWTSecret            string           `json:"jwtSecret"`
	SQLiteFile           string           `json:"sqliteFile"`
	Mail                 EmailCredentials `json:"mail"`
	OrderRetrievers      []*mail.Address  `json:"orderRetrievers"`
	RestMirror           RestMirror       `json:"restMirror"`
	TemplateInterval     uint             `json:"templateInterval"`
}

//Struct for all dav related paths
type DavPaths struct {
	PayedBills   string `json:"payedBills"`
	UnpayedBills string `json:"unpayedBills"`
	Scores       string `json:"scores"`
}

//Struct for SMTP credentials which will be used for sending mails.
type EmailCredentials struct {
	Username string        `json:"username"`
	Password string        `json:"password"`
	SMTPHost string        `json:"smtpHost"`
	Address  *mail.Address `json:"address"`
}

//Returns a default configuration with a generated jwt secret.
func defaultConf() *Configuration {
	secret := make([]byte, 16)
	rand.Read(secret)
	return &Configuration{
		Authentication: false,
		DavPaths: DavPaths{
			PayedBills:   "Rechnungen/bezahlt",
			UnpayedBills: "Rechnungen/offen",
			Scores:       "Noten",
		},
		GeneratedDirectory: "gen",
		Host:               "0.0.0.0:9400",
		Mail: EmailCredentials{
			Username: "max",
			Password: "thatsmaxmaildeliverypassword",
			SMTPHost: "mail.example.org",
			Address: &mail.Address{
				Name:    "Max Mustermann",
				Address: "noreply@mail.example.org",
			},
		},
		JWTExpirationMinutes: 5,
		JWTSecret:            fmt.Sprintf("%x", secret),
		SQLiteFile:           "core.sqlite",
		TemplateInterval:     5,
		RestMirror: RestMirror{
			Interval:           5,
			CategoriesPath:     "categories.json",
			CategoriesFlatPath: "categoriesflat.json",
			ScoresPath:         "scores.json",
		},
		OrderRetrievers: []*mail.Address{
			{
				Name:    "Max Mustermann",
				Address: "max.mustermann@mail.example.org",
			}}}
}

//section for the rest mirror
type RestMirror struct {
	Interval           uint   `json:"interval"`
	CategoriesPath     string `json:"categoriesPath"`
	CategoriesFlatPath string `json:"categoriesFlatPath"`
	ScoresPath         string `json:"scoresPath"`
}

//Reads configuration from file and creates one if it does not exist yet.
func config() (*Configuration, error) {
	conf, err := os.Open(configFile)
	if err != nil {
		conf, err = os.Create(configFile)
		if err != nil {
			return nil, err
		}
		encoder := json.NewEncoder(conf)
		encoder.SetIndent("", "  ")
		configuration := defaultConf()
		err = encoder.Encode(configuration)
		if err != nil {
			return configuration, err
		}
		return configuration, conf.Close()
	}
	decoder := json.NewDecoder(conf)
	configuration := &Configuration{}
	err = decoder.Decode(configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, conf.Close()
}
