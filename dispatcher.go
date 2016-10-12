package clerk

import (
	"fmt"
	"strings"
)

func HandleEmailApprove(filename string, send bool, c *Config) error {
	m, err := createApprovalEmail(filename, c)
	if err != nil {
		return err
	}

	timestamp, err := dateFromMDFilename(filename)
	if err != nil {
		return err
	}

	m.RenderMustache(map[string]string{"date": timestamp})

	if send {
		sc := ServerConfig{c.Mailserver.Hostname, c.Mailserver.Port}
		ap := new(AuthPair)
		ap.Read()
		if err := m.Send(&sc, ap); err != nil {
			return err
		}
	} else {
		m.DryRun(c.Author.Browser)
	}

	return nil
}

func HandleEmailToAll(filename string, send bool, c *Config) error {
	m, err := createAllEmail(filename, c)
	if err != nil {
		return err
	}

	timestamp, err := dateFromMDFilename(filename)
	if err != nil {
		return err
	}

	m.RenderMustache(map[string]string{"date": timestamp})

	if send {
		sc := ServerConfig{c.Mailserver.Hostname, c.Mailserver.Port}
		ap := new(AuthPair)
		ap.Read()
		if err := m.Send(&sc, ap); err != nil {
			return err
		}
	} else {
		m.DryRun(c.Author.Browser)
	}

	return nil
}

// Create an approval email
func createApprovalEmail(filename string, c *Config) (*PlainEmail, error) {
	html, errMD := MarkdownToHTML(filename)
	if errMD != nil {
		return nil, errMD
	}

	m := new(PlainEmail)
	m.FromName = c.Author.Name
	m.FromEmail = c.Author.Email
	m.ToEmails = c.Approval_list.Emails
	m.BccEmails = []string{c.Author.Email}
	m.Subject = c.Approval_list.Subject
	// TODO: Add mustache interpolation
	fmt.Printf("Will send to %v\n", m.ToEmails)
	bodyBegin := "<p>" + strings.Join([]string{
		c.Approval_list.Salutation,
		c.Approval_list.Text,
		c.Author.Notice,
	}, "</p>\r\n<p>") + "</p>\r\n\r\n"
	m.Body = append([]byte(bodyBegin), html...)

	return m, nil
}

func createAllEmail(filename string, c *Config) (*PlainEmail, error) {
	html, errMD := MarkdownToHTML(filename)
	if errMD != nil {
		return nil, errMD
	}

	m := new(PlainEmail)
	m.FromName = c.Author.Name
	m.FromEmail = c.Author.Email
	m.ToEmails = c.All_list.Emails
	m.BccEmails = []string{c.Author.Email}
	m.Subject = c.All_list.Subject
	// TODO: Add mustache interpolation
	fmt.Printf("Will send to %v\n", m.ToEmails)
	bodyBegin := "<p>" + strings.Join([]string{
		c.All_list.Salutation,
		c.All_list.Text,
		c.Author.Notice,
	}, "</p>\r\n<p>") + "</p>\r\n\r\n"
	m.Body = append([]byte(bodyBegin), html...)
	return m, nil
}

// func CreateIndividualEmails(filename string, c *Config) ([]PlainEmail, error) {
// }

// func filterNotifications () {}

func createSingleRecepientEmail(html []byte, r *Recipient, c *Config) (*PlainEmail, error) {
	m := new(PlainEmail)
	m.FromName = c.Author.Name
	m.FromEmail = c.Author.Email
	m.ToEmails = r.Emails
	m.BccEmails = []string{c.Author.Email}
	m.Subject = c.Individual_lists.Subject
	// TODO: Add mustache interpolation
	fmt.Printf("Will send to %v\n", m.ToEmails)
	bodyBegin := "<p>" + strings.Join([]string{
		r.Salutation,
		c.Individual_lists.Text,
		c.Author.Notice,
	}, "</p>\r\n<p>") + "</p>\r\n\r\n"
	m.Body = append([]byte(bodyBegin), html...)
	return m, nil
}
