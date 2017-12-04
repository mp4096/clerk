package clerk

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Author struct {
	Name    string
	Email   string
	Notice  string
	Browser string
}

type Recipients struct {
	Emails     []string
	Subject    string
	Salutation string
	Text       string
}

type EmailServer struct {
	Hostname string
	Port     int
}

type Config struct {
	Email_server    EmailServer
	Author          Author
	Distribute_list Recipients
	Approve_list    Recipients
}

// ImportFromFile reads configuration from a YAML file.
func (c *Config) ImportFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
