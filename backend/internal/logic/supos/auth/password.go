package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	return hashPassword(password)
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, password string) bool {
	return verifyPassword(hash, password)
}

func verifyPassword(hash, password string) bool {
	if hash == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
