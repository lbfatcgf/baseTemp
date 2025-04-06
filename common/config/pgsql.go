package config


type PgsqlConfig struct {
	Host        string  `yaml:"host     "`
	User        string  `yaml:"user     "`
	Password    string  `yaml:"password "`
	Dbname      string  `yaml:"dbname   "`
	Port        string  `yaml:"port     "`
	Sslmode     string  `yaml:"sslmode  "`
	TimeZone    string  `yaml:"TimeZone "`
	Sslrootcert string  `yaml:"sslrootcert"`
	Sslkey      string  `yaml:"sslkey   "`
	Sslcert     string  `yaml:"sslcert  "`
	Primary     bool    `yaml:"primary  "`
	Other       *string `yaml:"other"`
}

// GetHost 返回pgsql连接串
func (c *PgsqlConfig) GetHost() string {
	qulr := "host=" + c.Host + " " +
		"port=" + c.Port + " " +
		"user=" + c.User + " " +
		"password=" + c.Password + " " +
		"dbname=" + c.Dbname + " " +
		"sslmode=" + c.Sslmode + " " +
		"TimeZone=" + c.TimeZone
	if c.Sslmode != "disable" {
		qulr += " sslrootcert=" + c.Sslrootcert + " " +
			"sslkey=" + c.Sslkey + " " +
			"sslcert=" + c.Sslcert
	}
	if c.Other != nil && len(*c.Other) > 0 {
		qulr += " " + *c.Other
	}
	return qulr
}
