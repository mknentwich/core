package auth

import (
	"github.com/mknentwich/core/database"
)

func queryUserByEmail(email string) *database.User {
	var user database.User
	database.Receive().Where(&database.User{Email: email}).Find(&user)
	return &user
}

func saveUser(user *database.User) error {
	database.Receive().Save(user)
	if user.Password != "" {
		password, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		user := queryUserByEmail(user.Email)
		user.Password = password
		database.Receive().Save(user)
	}
	return nil
}
