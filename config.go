package clerk

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mailserver struct {
		Hostname string
		Port     int
	}
	Author struct {
		Name    string
		Email   string
		Notice  string
		Browser string
	}
	All_list struct {
		Emails     []string
		Subject    string
		Salutation string
		Text       string
	}
	Approval_list struct {
		Emails     []string
		Subject    string
		Salutation string
		Text       string
	}
	Individual_lists struct {
		Subject    string
		Text       string
		Recipients []Recipient
	}
}

type Recipient struct {
	ID         string
	Salutation string
	Emails     []string
}

func (c *Config) ImportFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	return nil
}
