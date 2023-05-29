package confutil

type Grpc struct {
	Etcd string `yaml:"etcd"`
	Port int    `yaml:"port"`
}

func (l *Grpc) GetHosts() string {
	return l.Etcd
}

func (l *Grpc) GetPort() int {
	return l.Port
}
