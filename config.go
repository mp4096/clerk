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
	EmailServer    EmailServer `yaml:"email_server"`
	Author         Author
	DistributeList Recipients `yaml:"distribute_list"`
	ApproveList    Recipients `yaml:"approve_list"`
}

// ImportFromFile reads configuration from a YAML file.
func (c *Config) ImportFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
