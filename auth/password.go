package auth

import (
	"github.com/mknentwich/core/database"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 8

//Hashes a given plaintext password. Returns the bcrypt's errors too.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	return string(hash), err
}

//Compares a plaintext password with the hashed one.
func comparePasswords(password string, user *database.User) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
