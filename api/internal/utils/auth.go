package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func PasswordMatchers(hashedPassword, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, InternalServerError{Message: err.Error()}
		}
	}

	return true, nil
}
