package auth

import (
	"github.com/mknentwich/core/database"
)

func queryUserByEmail(email string) *database.User {
	var user database.User
	database.Receive().Where(&database.User{Email: email}).Find(&user)
	return &user
}

func insertUser(user *database.User, password string) error {
	credentials := &Credentials{Email: user.Email, Password: password}
	database.Receive().Save(user)
	return updatePassword(credentials)
}

func updatePassword(credentials *Credentials) error {
	password, err := hashPassword(credentials.Password)
	if err != nil {
		return err
	}
	user := queryUserByEmail(credentials.Email)
	user.Password = password
	database.Receive().Save(user)
	return nil
}
