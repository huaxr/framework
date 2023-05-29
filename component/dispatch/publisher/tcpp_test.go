package publisher

import "testing"

func TestServer(t *testing.T) {
	initServerQ(8888)
}

func TestServer2(t *testing.T) {
	client()
}
