package config

import "net/url"

type RabbitMQConfig struct {
	Host     string `yaml:"host"    mapstructure:"host"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

func (c *RabbitMQConfig) GetLink() string {
	return "amqp://" + url.QueryEscape(c.User) + ":" + url.QueryEscape(c.Password) + "@" + c.Host
}
