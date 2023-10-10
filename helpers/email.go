package helpers

import (
	"bytes"
	"crypto/tls"
	"gin-user-management/initializers"
	"gin-user-management/models"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)


type EmailData struct {
	URL string
	FirstName string
	Subject string
}


func init() {
	initializers.LoadEnvVariables()
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}


func SendEmail(user *models.User, data *EmailData, emailTemplate string) {
	from := os.Getenv("EMAIL_FROM")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	to := user.Email
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	var body bytes.Buffer

	template, err := ParseTemplateDir("public/templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, emailTemplate, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	smtpPortInt, _ := ConvertStringToInt(smtpPort)
	d := gomail.NewDialer(smtpHost, smtpPortInt, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}


	// send email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}