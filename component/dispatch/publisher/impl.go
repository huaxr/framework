// Author: huaxr
// Time: 2022-12-08 15:48
// Git: huaxr

package publisher

import "io"

// Publisher multilateral support of Publisher, including nsq, kafka, tmp.
type Publisher interface {
	// Closer close
	io.Closer
	SetAuth(user string, passwd string)
	// Publish a message or event. extra of kafka: key + partition
	Publish(topic string, in io.Reader, extra ...interface{}) (err error)

	// setAuth before Run
	Run() error
}
