package encrypt

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(mainPassword, checkPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(mainPassword), []byte(checkPassword))

	return err == nil
}
