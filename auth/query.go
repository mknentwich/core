package auth

import (
	"github.com/mknentwich/core/database"
)

func queryUserByEmail(email string) *database.User {
	var user database.User
	err := database.Receive().Where(&database.User{UserWithoutPassword: &database.UserWithoutPassword{Email: email}}).Find(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func SaveUser(user *database.User) error {
	if us := queryUserByEmail(user.Email); us == nil {
		database.Receive().Save(user)
	} else {
		database.Receive().Model(user).Updates(user)
	}
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
