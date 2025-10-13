package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	genPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(genPassword), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
