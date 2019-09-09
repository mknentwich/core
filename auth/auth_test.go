package auth

import (
	"bytes"
	"encoding/json"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/utils"
	"io/ioutil"
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

//Appends Host URL.
func canonical(str string) string {
	return context.Conf.Host + "/auth" + str
}

//Compares both HTTP status and let the test fail if necessary.
func checkHttpStatus(expected int, status int, t *testing.T) {
	if expected != status {
		t.Errorf("HTTP status code should be %d but was %d!", expected, status)
	}
}

//Compares a credentials email with a users once and let the test fail if necessary.
func checkUserInfo(expected credentials, actual database.User, t *testing.T) {
	if expected.email != actual.Email {
		t.Errorf("Server returned wrong user. Expected email: %s, got %s", expected.email, actual.Email)
	}
}

//Sets the header for REST requests.
func header(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

//Tries to login and returns it response.
func login(credentials credentials, t *testing.T) *http.Response {
	data, err := json.Marshal(credentials)
	utils.Unexpected(t, err)
	request, err := http.NewRequest(http.MethodPost, canonical("/login"), bytes.NewBuffer(data))
	header(request)
	response, err := http.DefaultClient.Do(request)
	utils.Unexpected(t, err)
	return response
}

//Renews a given JWT and returns the new one and the HTTP status.
func renew(jwt string, t *testing.T) (string, int) {
	r, err := http.NewRequest(http.MethodGet, jwt, nil)
	r.Header.Set("Authorization: Bearer ", jwt)
	utils.Unexpected(t, err)
	response, err := http.DefaultClient.Do(r)
	utils.Unexpected(t, err)
	body, err := ioutil.ReadAll(response.Body)
	utils.Unexpected(t, err)
	return string(body), response.StatusCode
}

//Fetches a User from a JWT. Returns the HTTP status too.
func userInfo(jwt string, t *testing.T) (database.User, int) {
	r, err := http.NewRequest(http.MethodGet, canonical("/self"), nil)
	r.Header.Set("Authorization", "Bearer "+jwt)
	utils.Unexpected(t, err)
	response, err := http.DefaultClient.Do(r)
	utils.Unexpected(t, err)
	user := database.User{}
	json.NewDecoder(response.Body).Decode(&user)
	return user, response.StatusCode
}

//Tests a wrong password on `/login`
func TestInvalidCredentials(t *testing.T) {
	cr := albert
	cr.password += "sdklfgdklfj"
	response := login(cr, t)
	checkHttpStatus(http.StatusForbidden, response.StatusCode, t)
}

//Tests an invalid JWT format
func TestInvalidJWTFormat(t *testing.T) {
	_, status := userInfo("skdjfhskdjfhksjhrjkrg", t)
	checkHttpStatus(http.StatusUnprocessableEntity, status, t)
}

//Tests a JWT renewal with a valid JWT
func TestValidRenewal(t *testing.T) {
	response := login(albert, t)
	jwt, err := ioutil.ReadAll(response.Body)
	utils.Unexpected(t, err)
	renewJwt, status := renew(string(jwt), t)
	checkHttpStatus(http.StatusOK, status, t)
	user, status := userInfo(renewJwt, t)
	checkHttpStatus(http.StatusOK, status, t)
	checkUserInfo(albert, user, t)
}

//Tests the user info
func TestUserInfo(t *testing.T) {
	response := login(albert, t)
	jwt, err := ioutil.ReadAll(response.Body)
	utils.Unexpected(t, err)
	user, status := userInfo(string(jwt), t)
	checkHttpStatus(http.StatusOK, status, t)
	checkUserInfo(albert, user, t)
}
