package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/utils"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var albertUserWOP = database.UserWithoutPassword{Email: "albert@diealberts.at", Name: "Albert Albert", Admin: true}
var albertUser = database.User{
	UserWithoutPassword: &albertUserWOP,
	Password:            "desisagaunzgeheimespwd",
}
var albert = Credentials{Email: albertUser.Email, Password: albertUser.Password}
var config = context.Configuration{
	Authentication:       true,
	Host:                 "127.0.0.1:9400",
	JWTExpirationMinutes: 10,
	JWTSecret:            "sdvotkdoriuuuuuuuuuuuawbmet"}

func TestMain(m *testing.M) {
	go func() {
		err := context.InitializeCustomConfig(map[string]context.Serve{
			"db":   database.Serve,
			"auth": Serve},
			&config)
		if err != nil {
			fmt.Printf("Error during test setup: %s", err.Error())
			os.Exit(1)
		}
	}()
	r, err := http.NewRequest(http.MethodGet, canonical("/"), nil)
	if err != nil {
		fmt.Printf("Error during test setup: %s", err.Error())
		os.Exit(1)
	}
	c := http.Client{Timeout: time.Minute * 5}
	_, err = c.Do(r)
	for err != nil {
		time.Sleep(1 * time.Second)
		_, err = c.Do(r)
	}
	SaveUser(&albertUser)
	code := m.Run()
	os.Exit(code)
}

//Appends Host URL.
func canonical(str string) string {
	url := "http://" + config.Host + "/auth" + str
	fmt.Printf("Created URL: %s\n", url)
	return url
}

//Compares both HTTP status and let the test fail if necessary.
func checkHttpStatus(expected int, status int, t *testing.T) {
	if expected != status {
		t.Errorf("HTTP status code should be %d but was %d!", expected, status)
	}
}

//Compares a Credentials Email with a users once and let the test fail if necessary.
func checkUserInfo(expected Credentials, actual database.User, t *testing.T) {
	if expected.Email != actual.Email {
		t.Errorf("Server returned wrong user. Expected Email: %s, got %s", expected.Email, actual.Email)
	}
}

//Sets the header for REST requests.
func header(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

//Tries to login and returns it response.
func testLogin(credentials Credentials, t *testing.T) *http.Response {
	data, err := json.Marshal(credentials)
	utils.Unexpected(t, err)
	request, err := http.NewRequest(http.MethodPost, canonical("/login"), bytes.NewBuffer(data))
	utils.Unexpected(t, err)
	header(request)
	response, err := http.DefaultClient.Do(request)
	utils.Unexpected(t, err)
	return response
}

//Renews a given JWT and returns the new one and the HTTP status.
func renew(jwt string, t *testing.T) (string, int) {
	r, err := http.NewRequest(http.MethodGet, canonical("/refresh"), nil)
	utils.Unexpected(t, err)
	r.Header.Set("Authorization", "Bearer "+jwt)
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

func jwt(credentials Credentials) (string, bool) {
	data, _ := json.Marshal(credentials)
	request, err := http.NewRequest(http.MethodPost, canonical("/login"), bytes.NewBuffer(data))
	if err != nil {
		panic(err.Error())
	}
	header(request)
	response, _ := http.DefaultClient.Do(request)
	jwt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(jwt), response.StatusCode == http.StatusOK
}

func changePassword(issuerJwt string, user database.User) {
	data, _ := json.Marshal(user)
	request, err := http.NewRequest(http.MethodPost, canonical("/password"), bytes.NewBuffer(data))
	request.Header.Set("Authorization", "Bearer "+issuerJwt)
	if err != nil {
		panic(err.Error())
	}
	header(request)
	http.DefaultClient.Do(request)
}

//Tests a wrong Password on `/testLogin`
func TestInvalidCredentials(t *testing.T) {
	cr := albert
	cr.Password = "sdklfgdklfj"
	response := testLogin(cr, t)
	checkHttpStatus(http.StatusForbidden, response.StatusCode, t)
}

//Tests an invalid JWT format
func TestInvalidJWTFormat(t *testing.T) {
	_, status := userInfo("skdjfhskdjfhksjhrjkrg", t)
	checkHttpStatus(http.StatusUnprocessableEntity, status, t)
}

//Tests a JWT renewal with a valid JWT
func TestValidRenewal(t *testing.T) {
	response := testLogin(albert, t)
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
	response := testLogin(albert, t)
	jwt, err := ioutil.ReadAll(response.Body)
	utils.Unexpected(t, err)
	user, status := userInfo(string(jwt), t)
	checkHttpStatus(http.StatusOK, status, t)
	checkUserInfo(albert, user, t)
}

func TestUserUpdateWithoutPassword(t *testing.T) {
	albertName := "alberti"
	token, ok := jwt(albert)
	if !ok {
		panic("albert can't login")
	}
	nalbert := albertUser
	nalbert.Name = albertName
	nalbert.Password = ""
	changePassword(token, nalbert)
	changedAlbert := queryUserByEmail(nalbert.Email)
	if changedAlbert.Name != nalbert.Name {
		t.Errorf("albert's name should be %s but was %s", nalbert.Name, changedAlbert.Name)
	}
	if _, ok = jwt(albert); !ok {
		t.Errorf("albert cannot login anymore")
	}
	SaveUser(&albertUser)
}

func TestUserUpdatePassword(t *testing.T) {
	helgaWOP := database.UserWithoutPassword{
		Name:  "Helga",
		Email: "helga@gmx.at",
		Admin: false,
	}
	helga := database.User{
		UserWithoutPassword: &helgaWOP,
		Password:            "123456",
	}
	SaveUser(&helga)
	token, ok := jwt(Credentials{Email: helga.Email, Password: string(helga.Password)})
	if !ok {
		panic("helga cannot login")
	}
	helga.Password = "ibimsdehelga"
	changePassword(token, helga)
	if _, ok := jwt(Credentials{Email: helga.Email, Password: string(helga.Password)}); !ok {
		t.Errorf("helga cannot login with her new password")
	}
}

func TestAdminUpdateUser(t *testing.T) {
	williWOP := database.UserWithoutPassword{
		Name:  "Willi",
		Email: "willi@gmx.at",
		Admin: false}
	willi := database.User{
		UserWithoutPassword: &williWOP,
		Password:            "123456",
	}
	SaveUser(&willi)
	token, ok := jwt(albert)
	if !ok {
		panic("albert cannot login")
	}
	willi.Password = "ibimsdwilli"
	changePassword(token, willi)
	if _, ok := jwt(Credentials{Email: willi.Email, Password: string(willi.Password)}); !ok {
		t.Errorf("willi cannot login with his new password")
	}
}

func TestUserUpdateAnother(t *testing.T) {
	klausWOP := database.UserWithoutPassword{
		Name:  "Klaus",
		Email: "klaus@gmx.at",
		Admin: false}
	klaus := database.User{
		UserWithoutPassword: &klausWOP,
		Password:            "123456",
	}
	SaveUser(&klaus)
	token, ok := jwt(Credentials{Email: klaus.Email, Password: string(klaus.Password)})
	if !ok {
		panic("klaus cannot login")
	}
	nalbert := albertUser
	nalbert.Password = "se4vjktmhk"
	changePassword(token, nalbert)
	if _, ok := jwt(Credentials{Email: nalbert.Email, Password: string(nalbert.Password)}); ok {
		t.Errorf("klaus (who is no admin) was able to change another user's password")
	}
}
