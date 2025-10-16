package config

type AppConfig struct {
	Name                   string         `yaml:"name" mapstructure:"name"`
	Version                string         `yaml:"version" mapstructure:"version"`
	Mode                   string         `yaml:"mode" mapstructure:"mode"`
	Pgsql                  *[]PgsqlConfig    `yaml:"pgsql" mapstructure:"pgsql"`
	RabbitMQ               *[]RabbitMQConfig `yaml:"rabbitmq" mapstructure:"rabbitmq"`
	LogDir                 string         `yaml:"log_dir" mapstructure:"log_dir"`

}
