package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"strings"
	"time"
)

type jwtFormatError struct{}

func (*jwtFormatError) Error() string {
	return "invalid JWT structure"
}

//use for plainttext passwords only
type Credentials struct {
	Email    string
	Password string
}

//Defines the JWT header.
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

//Defines the JWT payload which is used to communicate between server and client.
type JWTPayload struct {
	Email      string    `json:"Email"`
	Expiration time.Time `json:"exp"`
	LastChange time.Time `json:"lastChange"`
}

//Does the testLogin process.
//Returns true, if the testLogin was successful.
//Returns the JWT for future request which require certain roles or permissions.
func login(data *Credentials) (string, error) {
	user := queryUserByEmail(data.Email)
	if user == nil {
		return "", errors.New("unknown user")
	}
	err := comparePasswords(data.Password, user)
	if err != nil {
		return "", err
	}
	token := generateToken(user)
	return token, nil
}

//Returns a JWT for a given Member.
func generateToken(user *database.User) string {
	header := JWTHeader{Algorithm: "HS256"}
	payload := JWTPayload{Email: user.Email, Expiration: time.Now().Add(time.Minute * time.Duration(context.Conf.JWTExpirationMinutes))}
	head, _ := json.Marshal(header)
	pay, _ := json.Marshal(payload)
	encodedHead := base64.URLEncoding.EncodeToString([]byte(head))
	encodedPayload := base64.URLEncoding.EncodeToString([]byte(pay))
	rawToken := encodedHead + "." + encodedPayload
	fullToken := rawToken + "." + hash(rawToken)
	return fullToken
}

//Hashes a JWT.
func hash(rawToken string) string {
	sig := hmac.New(sha256.New, []byte(context.Conf.JWTSecret))
	sig.Write([]byte(rawToken))
	return hex.EncodeToString(sig.Sum(nil))
}

//Checks, if the given JWT is valid.
//Returns true if it is valid.
//Returns the user which belongs to the JWT. If the JWT is not valid, the user will be nil.
func Check(token string) (user *database.User, expiration time.Time, err error) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 || hash(tokenParts[0]+"."+tokenParts[1]) != tokenParts[2] {
		err = &jwtFormatError{}
		return
	}
	payloadJsonByte, _ := base64.URLEncoding.DecodeString(tokenParts[1])
	payload := &JWTPayload{}
	json.Unmarshal(payloadJsonByte, payload)
	if time.Now().After(payload.Expiration) {
		err = errors.New("jwt expired")
		return
	}
	user = queryUserByEmail(payload.Email)
	expiration = payload.Expiration
	return
}
