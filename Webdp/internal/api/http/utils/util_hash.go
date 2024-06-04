package utils

import (
	"fmt"

	errors "webdp/internal/api/http"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrUnexpected, err.Error())
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteSavedHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteSavedHash, bytePlain)
	return err == nil
}
