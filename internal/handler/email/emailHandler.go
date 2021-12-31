package email

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/vinicius.csantos/nagini-api/config"
	"github.com/vinicius.csantos/nagini-api/internal/model"
	constantUtils "github.com/vinicius.csantos/nagini-api/internal/util/constant"
	"github.com/vinicius.csantos/nagini-api/internal/util/encrypt"
	"html/template"
	"net/smtp"
	"time"
)

var baseWebUrl = config.Config("WEB_URL", "localhost:3000")

func SendEmailToConfirmAccount(email, username, cpfCNPJ string, userID uuid.UUID) error {

	timeEncrypted := time.Now()

	stringToEncrypt := fmt.Sprintf("%s$%s$%s$%s", email, username, cpfCNPJ, timeEncrypted)

	token, err := createTokenAndSaveItInCache(stringToEncrypt, userID, timeEncrypted)

	if err != nil {
		return err
	}

	accountConfirmationUrl := fmt.Sprintf(constantUtils.GeneralUrlFormat, baseWebUrl, "confirm-account", token)

	templatePath := fmt.Sprintf("%s/confirm-email.html", config.Config("EMAIL_TEMPLATES_PATH", ""))

	confirmStruct := struct {
		Username    string
		ActivateUrl string
	}{
		Username:    username,
		ActivateUrl: accountConfirmationUrl,
	}

	subject := " Email Confirmation! \n%s\n\n"

	err = sendEmail(email, templatePath, subject, confirmStruct)

	return err
}

func SendEmailToRecoverPassword(email, username, cpfCNPJ string, userID uuid.UUID) error {

	timeEncrypted := time.Now()

	stringToEncrypt := fmt.Sprintf("%s$%s$%s$%s", email, username, cpfCNPJ, timeEncrypted)

	token, err := createTokenAndSaveItInCache(stringToEncrypt, userID, timeEncrypted)

	if err != nil {
		return err
	}

	passwordRecoverUrl := fmt.Sprintf(constantUtils.GeneralUrlFormat, baseWebUrl, "change-password", token)

	recoverCancelUrl := fmt.Sprintf(constantUtils.GeneralUrlFormat, baseWebUrl, "cancel-recover-password", token)

	templatePath := fmt.Sprintf("%s/recover-password.html", config.Config("EMAIL_TEMPLATES_PATH", ""))

	myStruct := struct {
		Username   string
		RecoverUrl string
		CancelLink string
	}{
		Username:   username,
		RecoverUrl: passwordRecoverUrl,
		CancelLink: recoverCancelUrl,
	}

	if err != nil {
		return err
	}

	subject := " Password Recovery! \n%s\n\n"

	err = sendEmail(email, templatePath, subject, myStruct)

	return err
}

func createTokenAndSaveItInCache(stringToEncrypt string, userID uuid.UUID, timeToEncrypt time.Time) (string, error) {
	userCache := model.UserCache
	token, err := encrypt.UrlEncrypt(stringToEncrypt)
	cacheKey := fmt.Sprintf("%s-%s", userID.String(), timeToEncrypt.String())

	if err != nil {
		return "", err
	}
	err = userCache.SetTTL(30 * time.Minute)

	if err != nil {
		return "", err
	}

	err = userCache.Set(cacheKey, token)

	if err != nil {
		return "", err
	}

	return token, nil
}

func sendEmail(email, templatePath, subject string, structure interface{}) error {

	from := config.Config("EMAIL_SEND", "")
	password := config.Config("EMAIL_PASSWORD", "")

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	emailTemplate, _ := template.ParseFiles(templatePath)

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body.Write([]byte(fmt.Sprintf("Subject:"+subject, mimeHeaders)))

	err := emailTemplate.Execute(&body, structure)

	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())

	return err
}
