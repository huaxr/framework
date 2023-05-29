package confutil

import (
	"testing"
	"time"
)

type T struct {
	// you must set A instead of AA cause yml bugs
	AA string `yaml:"a"`
	A  string `yaml:"a"`
}

var x = new(T)

func TestConf(t *testing.T) {

	time.Now()
	GetDefaultConfig()
	InitYmlFile("test.yml", x)
	t.Log(x.AA)
	t.Log(x.A)
}
