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
		builder = builder.AddRecipients(&c.ApproveList)
	case DISTRIBUTE:
		builder = builder.AddRecipients(&c.DistributeList)
	default:
		return errors.New("Unknown action type.")
	}
	email, err := builder.Build(context)
	if err != nil {
		return err
	}

	if send {
		fmt.Printf("Will send to %v\n", email.GetRecipients())
		fmt.Printf("Please enter your credentials for \"%s\"\n", c.EmailServer.Hostname)
		ap := new(authPair)
		ap.prompt()
		if err := email.Send(&c.EmailServer, ap); err != nil {
			return err
		}
	} else {
		fmt.Printf("Send flag not set: opening preview in \"%s\"\n", c.Author.Browser)
		if err := email.OpenInBrowser(c.Author.Browser); err != nil {
			return err
		}
	}

	return nil
}
