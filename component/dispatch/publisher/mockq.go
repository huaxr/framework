// Author: huaxr
// Time:   2022/1/4 下午2:57
// Git:    huaxr

package publisher

import (
	"context"
	"io"
	"io/ioutil"
)

type mockPublisher struct {
	topicChan map[string]chan []byte
}

func (p *mockPublisher) Publish(topic string, in io.Reader, infors ...interface{}) (err error) {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	p.topicChan[topic] <- body
	return nil
}

func (p *mockPublisher) Run() error { return nil }

func (p *mockPublisher) SetAuth(u string, pa string) {}

func (p *mockPublisher) Close() error { return nil }

func NewMOCKPublisher(ctx context.Context) Publisher {
	return nil
}
