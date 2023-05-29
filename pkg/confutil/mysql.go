package confutil

type Mysql struct {
	Slaves       []string `yaml:"slaves"`
	Master       string   `yaml:"master"`
	MaxConn      int      `yaml:"maxConn"`
	MaxIdle      int      `yaml:"maxIdle"`
	ShowSql      bool     `yaml:"showSql"`
	SlowDuration int      `yaml:"slowDuration"`
}
