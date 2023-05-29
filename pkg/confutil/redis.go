package confutil

type Redis struct {
	Host        string `yaml:"host"`
	Password    string `yaml:"password"`
	Idletimeout int    `yaml:"idletimeout"`
	Readtimeout int    `yaml:"readtimeout"`
	MaxRetry    int    `yaml:"maxretry"`
	Poolsize    int    `yaml:"poolsize"`
	Db          int    `yaml:"db"`
	ShowQuery   bool   `yaml:"showQuery"`
}
