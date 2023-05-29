// Author: huaxr
// Time:   2022/6/9 下午3:50
// Git:    huaxr

package publisher

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"

	"github.com/huaxr/framework/logx"
)

var eof = []byte{0, 1, 0, 1}

type conn struct {
	conn net.Conn
}

type tcpPublisher struct {
	l net.Listener
}

func (t *tcpPublisher) SetAuth(u string, pa string) {}

func (t *tcpPublisher) Run() error { return nil }

func (t *tcpPublisher) Publish(topic string, in io.Reader, extra ...interface{}) (err error) {
	return nil
}

func (t *tcpPublisher) Close() error { return nil }

func do(b []byte) {
	logx.L().Debugf(string(b))
}

func resp() []byte {
	return append([]byte("123"), eof...)
}

func handler(c *conn, f func([]byte)) {
	defer func() {
		_ = c.conn.Close()
	}()

	var req = make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			return
		}
		req = append(req, buffer[:n]...)
		if len(req) >= len(eof) && reflect.DeepEqual(req[len(req)-len(eof):], eof) {
			f(req[:len(req)-len(eof)])

			_, _ = c.conn.Write(resp())

			req = req[:0]
		}
	}
}

func initServerQ(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logx.L().Errorf("err tcp setup, %v", err)
		return
	}

	tcp := &tcpPublisher{
		l: listener,
	}

	for {
		con, err := tcp.l.Accept()
		if err != nil {
			log.Println("accept err", err)
			continue
		}

		c := &conn{
			conn: con,
		}
		go handler(c, do)
	}
}

func identify() []byte {
	return append([]byte("hello server"), eof...)
}

func client() {
	c, err := net.Dial("tcp", "0.0.0.0:8888")
	if err != nil {
		return
	}

	conn := &conn{conn: c}

	// first send a byte, we can identify client info here.
	_, _ = conn.conn.Write(identify())

	handler(conn, do)
}
