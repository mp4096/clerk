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

type authPair struct {
	login    string
	password string
}

// Asks user for the login and password for the email server
func (ap *authPair) prompt() error {
	fmt.Printf("Login: ")
	fmt.Scanln(&ap.login)

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	ap.password = string(pass) // `pass` is already a slice

	return nil
}

type EmailBuilder interface {
	AddAuthor(*Author) EmailBuilder
	AddRecipients(*Recipients) EmailBuilder
	AddContent(string) EmailBuilder
	FillInContext(map[string]string) EmailBuilder
	Build() Email
}

// Email builder struct containing all the nitty-gritty details and subfields of our custom email
type emailBuilder struct {
	fromName        string
	fromEmail       string
	recipientEmails []string
	bccEmails       []string
	subject         string
	salutation      string
	headerText      string
	mailText        string
	notice          string
}

func NewEmail() EmailBuilder {
	return &emailBuilder{}
}

func (eb *emailBuilder) AddAuthor(a *Author) EmailBuilder {
	eb.fromName = a.Name
	eb.fromEmail = a.Email
	eb.bccEmails = []string{a.Email}
	eb.notice = a.Notice

	return eb
}

func (eb *emailBuilder) AddRecipients(rs *Recipients) EmailBuilder {
	eb.recipientEmails = rs.Emails
	eb.subject = rs.Subject
	eb.salutation = rs.Salutation
	eb.headerText = rs.Text

	return eb
}

func (eb *emailBuilder) AddContent(s string) EmailBuilder {
	eb.mailText = s

	return eb
}

func (eb *emailBuilder) FillInContext(context map[string]string) EmailBuilder {
	eb.subject = mustache.Render(eb.subject, context)
	eb.headerText = mustache.Render(eb.headerText, context)
	eb.mailText = mustache.Render(eb.mailText, context)

	return eb
}

func (eb *emailBuilder) Build() Email {
	body := "From: " + eb.fromName + "<" + eb.fromEmail + ">\r\n"
	body += "To: " + strings.Join(eb.recipientEmails, ", ") + "\r\n"
	body += "Subject: " + eb.subject + "\r\n"
	body += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body += "<p>" + eb.salutation + "</p>\r\n"
	body += "<p>" + eb.headerText + "</p>\r\n"
	body += "<p>" + eb.notice + "</p>\r\n\r\n"
	body += eb.mailText

	e := new(email)
	e.fromEmail = eb.fromEmail
	e.toEmails = append(eb.recipientEmails, eb.bccEmails...)
	e.body = []byte(body)

	return e
}

type Email interface {
	Send(*EmailServer, *authPair) error
	OpenInBrowser(string) error
	GetRecipients() []string
}

type email struct {
	fromEmail string
	toEmails  []string
	body      []byte
}

func (e *email) Send(s *EmailServer, ap *authPair) error {
	serverInfo := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	auth := smtp.PlainAuth("", ap.login, ap.password, s.Hostname)

	return smtp.SendMail(serverInfo, auth, e.fromEmail, e.toEmails, e.body)
}

func (e *email) OpenInBrowser(browserName string) error {
	html := append(
		[]byte(`<html><head><meta charset="UTF-8"></head>`),
		e.body...,
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

	cmd := exec.Command(browserName, tmpfile.Name())
	return cmd.Run()
}

func (e *email) GetRecipients() []string {
	return e.toEmails
}
