// Author: huaxr
// Time:   2021/8/6 下午2:12
// Git:    huaxr

package consumer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/huaxr/framework/logx"

	"github.com/nsqio/go-nsq"
)

var (
	nsqLogger = log.New(ioutil.Discard, "", log.LstdFlags)
)

type nsqConfig struct {
	token       string
	namespace   string
	maxInflight int

	nsqAddr    []string
	lookupAddr []string
	concurrent int
}

type nsqConsumer struct {
	ctx    context.Context
	topic  string
	config *nsqConfig
	slave  *unsafe.Pointer
}

func NewNSQConsumer(ctx context.Context, topic string, hosts []string, lookups []string, concurrent int) Consumer {
	return &nsqConsumer{
		ctx:   ctx,
		topic: topic,
		config: &nsqConfig{
			nsqAddr:    hosts,
			lookupAddr: lookups,
			concurrent: concurrent,
		},
	}
}

func (h *nsqConsumer) SetAuth(u string, p string) {}

// HandleMessage is implement of nsq.Consume.Handler.
// which handle the nsq.message and impetus the in-flight loop works.
// when return error != nil, OnRequeue will be triggered. else OnFinish
// will tell the nsqd while !m.DisableAutoResponse()
func (h *nsqConsumer) HandleMessage(m *nsq.Message) (err error) {
	logx.L().Debugf(string(m.Body))
	return
}

// InitConsumer nsq worker as a partitionConsumer and send response to server. server become a distribute center.
func (h *nsqConsumer) Run() error {
	config := nsq.NewConfig()
	// the mq cluster configure the auth-url with tmp-env address link.
	// the /auth api will never reached in local environment.
	config.AuthSecret = fmt.Sprintf("%v?%v", h.config.token, h.config.namespace)

	if h.config.maxInflight != 0 {
		config.MaxInFlight = h.config.maxInflight
	}
	config.Snappy = true
	consumer, err := nsq.NewConsumer(h.topic, "nil", config)
	if err != nil {
		logx.L().Errorf("nsq err:%v", err)
		return err
	}
	consumer.SetLogger(nsqLogger, nsq.LogLevelDebug)
	// partitionConsumer -> &partitionConsumer will panic.
	// the parameter of unsafe.Pointer is pointer already.
	h.setConsumer(unsafe.Pointer(consumer))
	return nil
}

func (h *nsqConsumer) connectDirectly() {
	err := h.getConsumer().ConnectToNSQDs(h.config.nsqAddr)
	if err != nil {
		panic(err)
	}
	stats := h.getConsumer().Stats()
	if stats.Connections == 0 {
		panic("stats report 0 connections (should be > 0)")
	}
}

func (h *nsqConsumer) connectLookUp() {
	c := h.getConsumer()
	if err := c.ConnectToNSQLookupds(h.config.lookupAddr); err != nil {
		log.Printf("connectLookUp err %v", err)
		panic(err)
	}
}

func (h *nsqConsumer) Consume() {
	h.getConsumer().AddConcurrentHandlers(h, h.config.concurrent)
	h.connectLookUp()
	logx.L().Debugf("start %d nsq consumers", h.config.concurrent)
}

func (h *nsqConsumer) Close() {
	h.getConsumer().Stop()
}

func (h *nsqConsumer) getConsumer() *nsq.Consumer {
	return (*nsq.Consumer)(*h.slave)
}

func (h *nsqConsumer) setConsumer(pointer unsafe.Pointer) {
	h.slave = &pointer
}
