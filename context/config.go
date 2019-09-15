package context

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const configFile = "config.json"
const customConfigFile = "customConfig.json"

var Conf Configuration

//Struct for the configuration of the application.
type Configuration struct {
	Authentication       bool   `json:"authentication"`
	Host                 string `json:"host"`
	JWTExpirationMinutes int    `json:"jwtExpirationMinutes"`
	JWTSecret            string `json:"jwtSecret"`
	SQLiteFile           string `json:"sqliteFile"`
}

//Returns a default configuration with a generated jwt secret.
func defaultConf() *Configuration {
	secret := make([]byte, 16)
	rand.Read(secret)
	return &Configuration{
		Authentication:       true,
		Host:                 "0.0.0.0:9400",
		JWTExpirationMinutes: 5,
		JWTSecret:            fmt.Sprintf("%x", secret),
		SQLiteFile:           "core.sqlite"}
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
	return configuration, err
	if err != nil {
		return configuration, err
	}
	return configuration, conf.Close()
}
