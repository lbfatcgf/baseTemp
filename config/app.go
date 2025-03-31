package config

type AppConfig struct {
	Name                   string         `yaml:"name" mapstructure:"name"`
	Version                string         `yaml:"version" mapstructure:"version"`
	Mode                   string         `yaml:"mode" mapstructure:"mode"`
	Pgsql                  *PgsqlConfig    `yaml:"pgsql" mapstructure:"pgsql"`
	RabbitMQ               *RabbitMQConfig `yaml:"rabbitmq" mapstructure:"rabbitmq"`
	LogDir                 string         `yaml:"log_dir" mapstructure:"log_dir"`
	InitTable              bool           `yaml:"init_table" mapstructure:"init_table"`
	SendToMq               bool           `yaml:"send_to_mq" mapstructure:"send_to_mq"`
}
