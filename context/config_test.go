package context

import (
	"encoding/json"
	"os"
	"testing"
)

//Fails a test if an unexpected error occur which has nothing to do with the test case.
func unexpected(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Encountered unexpected error: %s", err.Error())
	}
}

//Removes config file before and after every test case.
func TestMain(m *testing.M) {
	os.Remove(configFile)
	code := m.Run()
	os.Remove(configFile)
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
	unexpected(t, err)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(custom)
	unexpected(t, err)
	err = file.Close()
	unexpected(t, err)
	read, err := config()
	unexpected(t, err)
	if read.Host != custom.Host || read.JWTExpirationMinutes != custom.JWTExpirationMinutes || read.SQLiteFile != custom.SQLiteFile {
		t.Errorf("Expected configuration: %v, but got: %v", *custom, *read)
	}
}
