package clerk

import (
	"errors"
	"fmt"
)

type Action uint8

const (
	NOTHING Action = iota
	APPROVE
	DISTRIBUTE
)

func ProcessFile(filename string, a Action, send bool, c *Config) error {
	timestamp, err := dateFromMDFilename(filename)
	if err != nil {
		return err
	}
	context := map[string]string{"date": timestamp}

	md, err := markdownToHTML(filename)
	if err != nil {
		return err
	}

	builder := NewEmail().AddAuthor(&c.Author).AddContent(md)
	switch a {
	case APPROVE:
		builder = builder.AddRecipients(&c.Approve_list)
	case DISTRIBUTE:
		builder = builder.AddRecipients(&c.Distribute_list)
	default:
		return errors.New("Unknown action type.")
	}
	email := builder.FillInContext(context).Build()

	if send {
		fmt.Printf("Will send to %v\n", email.GetRecipients())
		ap := new(authPair)
		ap.prompt()
		if err := email.Send(&c.Email_server, ap); err != nil {
			return err
		}
	} else {
		fmt.Println("Opening preview in browser")
		if err := email.OpenInBrowser(c.Author.Browser); err != nil {
			return err
		}
	}

	return nil
}
