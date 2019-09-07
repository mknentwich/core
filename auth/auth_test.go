package auth

import (
	"bytes"
	"encoding/json"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/utils"
	"net/http"
	"os"
	"testing"
)

var albert = credentials{email: "albert@diealberts.at", password: "thisisaverysecretpassphraseandshouldnotbepublished"}

func TestMain(m *testing.M) {
	context.InitializeCustomConfig(map[string]context.Serve{
		"/db":   database.Serve,
		"/auth": Serve},
		&context.Configuration{
			Host:      "http://127.0.0.1:9400",
			JWTSecret: "sdvotkdoriuuuuuuuuuuuawbmet"})
	code := m.Run()
	os.Exit(code)
}

func call(str string) string {
	return context.Conf.Host + "/auth" + str
}

func header(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

func post(body interface{}, t *testing.T) *http.Response {
	data, err := json.Marshal(body)
	utils.Unexpected(t, err)
	request, err := http.NewRequest(http.MethodPost, call("/login"), bytes.NewBuffer(data))
	header(request)
	response, err := http.DefaultClient.Do(request)
	utils.Unexpected(t, err)
	return response
}

func TestInvalidCredentials(t *testing.T) {
	cr := albert
	cr.password += "sdklfgdklfj"
	response := post(cr, t)
	status := http.StatusForbidden
	if response.StatusCode != status {
		t.Errorf("HTTP status code should be %d but was %d!", status, response.StatusCode)
	}
}

func TestInvalidRenewal(t *testing.T) {}

func TestValidRenewal(t *testing.T) {}

func TestValidUserInfo(t *testing.T) {}

func TestInvalidUserInfo(t *testing.T) {}
