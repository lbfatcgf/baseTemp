package config

import "net/url"

type RabbitMQConfig struct {
	Host     string `yaml:"host"    mapstructure:"host"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

func (c *RabbitMQConfig) GetLink() string {
	return url.QueryEscape("amqp://" + c.User + ":" + c.Password + "@" + c.Host)
}
