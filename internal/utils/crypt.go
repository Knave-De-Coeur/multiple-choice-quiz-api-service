package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt uses bcrypt package to encrypt pass
func HashAndSalt(pwd []byte) (string, error) {
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePasswords uses bcyrpt library to check the stored password with the plain string password
func ComparePasswords(hashedPwd string, plainPwd []byte) (bool, error) {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
