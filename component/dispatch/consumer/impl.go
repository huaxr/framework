// Author: huaxr
// Time: 2022-12-08 15:48
// Git: huaxr

package consumer

type Consumer interface {
	// SetAuth before InitConsumer
	SetAuth(user string, passwd string)
	Run() error

	// Run then after Consume
	Consume()
}
