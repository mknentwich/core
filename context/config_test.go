package context

import (
	"encoding/json"
	"github.com/mknentwich/core/utils"
	"os"
	"testing"
)

//Removes config file before and after every test case.
func TestMain(m *testing.M) {
	os.Remove(configFile)
	os.Remove(customConfigFile)
	code := m.Run()
	os.Remove(configFile)
	os.Remove(customConfigFile)
	os.Exit(code)
}

//Tests if a config will be generated when no exists.
func TestConfigCreation(t *testing.T) {
	_, err := config()
	if err != nil {
		t.Errorf("Error on creating a config file: %s", err.Error())
	}
}

//Tests if the read config is the persisted one.
func TestConfigRead(t *testing.T) {
	custom := &Configuration{Host: "73.73.73.73:80", JWTExpirationMinutes: 1, SQLiteFile: "db.sqlite"}
	file, err := os.Create(configFile)
	utils.Unexpected(t, err)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(custom)
	utils.Unexpected(t, err)
	err = file.Close()
	utils.Unexpected(t, err)
	read, err := config()
	utils.Unexpected(t, err)
	if read.Host != custom.Host || read.JWTExpirationMinutes != custom.JWTExpirationMinutes || read.SQLiteFile != custom.SQLiteFile {
		t.Errorf("Expected configuration: %v, but got: %v", *custom, *read)
	}
}
