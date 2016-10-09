package clerk

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os/exec"
	"strings"

	"github.com/hoisie/mustache"
	"github.com/howeyc/gopass"
)

type PlainEmail struct {
	FromName  string
	FromEmail string
	ToEmails  []string
	Subject   string
	Body      []byte
}

type AuthPair struct {
	Login    string
	Password string
}

func (ap *AuthPair) Read() error {
	fmt.Printf("Login: ")
	fmt.Scanln(&ap.Login)

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return err
	}

	ap.Password = string(pass) // TODO: is this the best way to convert bytes to string?

	return nil
}

type ServerConfig struct {
	Hostname string
	Port     int
}

func (m *PlainEmail) RenderMustache(context map[string]string) {
	m.Subject = string(mustache.Render(string(m.Subject), context))
	m.Body = []byte(mustache.Render(string(m.Body), context))
}

func (m *PlainEmail) ExportHTML() []byte {
	fieldFrom := "From: " + m.FromName + "<" + m.FromEmail + ">\r\n"
	fieldTo := "To: " + strings.Join(m.ToEmails, ", ") + "\r\n"
	fieldSubject := "Subject: " + m.Subject + "\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	return append([]byte(fieldFrom+fieldTo+fieldSubject+mime), m.Body...)
}

func (m *PlainEmail) Send(sc *ServerConfig, ap *AuthPair) error {
	server := fmt.Sprintf("%s:%d", sc.Hostname, sc.Port)
	auth := smtp.PlainAuth("", ap.Login, ap.Password, sc.Hostname)
	err := smtp.SendMail(server, auth, m.FromEmail, m.ToEmails, m.ExportHTML())
	if err != nil {
		return err
	}
	return nil
}

func (m *PlainEmail) DryRun(browser string) error {
	html := append(
		[]byte(`<html><head><meta charset="UTF-8"></head>`),
		m.ExportHTML()...,
	)
	html = append(html, []byte(`</html>`)...)

	tmpfile, err := ioutil.TempFile(".", "clerk_preview_")
	if err != nil {
		return err
	}
	if _, err := tmpfile.Write(html); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := exec.Command(browser, tmpfile.Name())
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
