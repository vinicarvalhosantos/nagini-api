package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/vinicius.csantos/nagini-api/config"
	"golang.org/x/crypto/bcrypt"
	"net/url"
)

var secretToken = config.GetSecretKey("SECRET_BASE64_KEY")
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(mainPassword, checkPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(mainPassword), []byte(checkPassword))

	return err == nil
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func UrlEncrypt(textToEncrpyt string) (string, error) {
	block, err := aes.NewCipher([]byte(secretToken))
	if err != nil {
		return "", err
	}

	plainText := []byte(textToEncrpyt)
	cfbEncrypter := cipher.NewCFBEncrypter(block, bytes)

	encrypted := make([]byte, len(plainText))
	cfbEncrypter.XORKeyStream(encrypted, plainText)

	return url.QueryEscape(encode(encrypted)), nil
}

func UrlDecrypt(str string) (string, error) {
	block, err := aes.NewCipher([]byte(secretToken))
	if err != nil {
		return "", nil
	}

	str, err = url.QueryUnescape(str)

	if err != nil {
		return "", err
	}

	decrypted, err := decode(str)

	if err != nil {
		return "", err
	}

	cfbDecrypter := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(decrypted))
	cfbDecrypter.XORKeyStream(plainText, decrypted)

	return string(plainText), nil
}
