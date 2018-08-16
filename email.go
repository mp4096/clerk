package clerk

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
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
	Build(map[string]string) (Email, error)
}

// Email builder struct containing all the nitty-gritty details and subfields of our custom email
type emailBuilder struct {
	fromName        string
	fromEmail       string
	recipientEmails []string
	bccEmails       []string
	subject         string
	salutation      string
	preludeText     string
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
	eb.preludeText = rs.Text

	return eb
}

func (eb *emailBuilder) AddContent(s string) EmailBuilder {
	eb.mailText = s

	return eb
}

func (eb *emailBuilder) Build(context map[string]string) (Email, error) {
	mustache.AllowMissingVariables = false

	renderedSubject, err := mustache.Render(eb.subject, context)
	if err != nil {
		return nil, err
	}
	headerTemplate := "From: " + encodeRfc1342(eb.fromName) + " <" + eb.fromEmail + ">\r\n"
	headerTemplate += "To: " + strings.Join(eb.recipientEmails, ", ") + "\r\n"
	headerTemplate += "Subject: " + encodeRfc1342(renderedSubject) + "\r\n"
	headerTemplate += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	header, err := mustache.Render(headerTemplate, context)
	if err != nil {
		return nil, err
	}

	bodyTemplate := "<p>" + eb.salutation + "</p>\r\n"
	bodyTemplate += "<p>" + eb.preludeText + "</p>\r\n"
	bodyTemplate += "<p>" + eb.notice + "</p>\r\n\r\n"
	bodyTemplate += eb.mailText
	body, err := mustache.Render(bodyTemplate, context)
	if err != nil {
		return nil, err
	}

	e := new(email)
	e.fromEmail = eb.fromEmail
	e.toEmails = append(eb.recipientEmails, eb.bccEmails...)
	e.header = []byte(header)
	e.body = []byte(body)

	return e, nil
}

type Email interface {
	Send(*EmailServer, *authPair) error
	OpenInBrowser(string) error
	GetRecipients() []string
}

type email struct {
	fromEmail string
	toEmails  []string
	header    []byte
	body      []byte
}

func (e *email) Send(s *EmailServer, ap *authPair) error {
	serverInfo := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	auth := smtp.PlainAuth("", ap.login, ap.password, s.Hostname)

	return smtp.SendMail(serverInfo, auth, e.fromEmail, e.toEmails, append(e.header, e.body...))
}

func (e *email) OpenInBrowser(browserName string) error {
	html := append(
		[]byte("<html><head><meta charset=\"UTF-8\"></head>\n<pre>"),
		escapeAngleBrackets(e.header)...,
	)
	html = append(html, []byte("</pre>\n")...)
	html = append(html, e.body...)
	html = append(html, []byte("</html>")...)

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

	oldFilename, err := filepath.Abs(filepath.Join(".", tmpfile.Name()))
	if err != nil {
		return err
	}
	newFilename := oldFilename + ".html"
	os.Rename(oldFilename, newFilename)

	cmd := exec.Command(browserName, newFilename)
	return cmd.Start()
}

func (e *email) GetRecipients() []string {
	return e.toEmails
}
