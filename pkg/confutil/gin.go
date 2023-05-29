package confutil

type mode string

const (
	Debug   mode = "debug"
	Release mode = "release"
	Tester  mode = "test"
)

func (m mode) String() string {
	return string(m)
}

type Gin struct {
	Mode mode `yaml:"mode"`
	Port int  `yaml:"port"`
}

func (l *Gin) GetPort() int {
	return l.Port
}
