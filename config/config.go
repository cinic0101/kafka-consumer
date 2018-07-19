package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Kafka struct {
		BootstrapServers string `yaml:"bootstrap-servers"`
		GroupID string `yaml:"group-id"`
		Topics string `yaml:"topics"`
	}
	Destination struct {
		Type string
		Params map[string]string
	}
}

func Unmarshal(filename string) *Config {
	c := &Config{}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(buf, &c); err != nil {
		panic(err)
	}

	return c
}