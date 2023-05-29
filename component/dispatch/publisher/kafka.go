// Author: huaxr
// Time:   2022/1/26 上午11:41
// Git:    huaxr

package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/internal/define"

	"github.com/huaxr/framework/logx"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

type partitionWay string

const (
	Hash   partitionWay = "hash"
	Manual partitionWay = "manual"
	Random partitionWay = "random"
)

var (
	kafkaProducer *kafkaPublisher
)

type kafkaPublisher struct {
	define.Opening
	sync.RWMutex
	ctx          context.Context
	user, passwd string
	// binding address with client
	producer         sarama.AsyncProducer
	hosts            []string
	partition        partitionWay
	failFuncCallback func(*sarama.ProducerMessage, error)
	close            chan struct{}
}

// signal calling
func NewKafkaPublisher(ctx context.Context, hosts []string, partition partitionWay, errCallback func(*sarama.ProducerMessage, error)) Publisher {
	kafkaProducer = new(kafkaPublisher)
	kafkaProducer.hosts = hosts
	kafkaProducer.ctx, _ = context.WithCancel(ctx)
	kafkaProducer.partition = partition
	kafkaProducer.failFuncCallback = errCallback
	kafkaProducer.close = make(chan struct{})
	return kafkaProducer
}

func (p *kafkaPublisher) Close() error {
	logx.L().Debugf("close kafka connection")
	err := p.producer.Close()
	if err != nil {
		logx.T(nil, define.ArchError).Errorf("kafkaPublisher Close err:%v", err)
	}
	p.close <- struct{}{}
	return err
}

func (p *kafkaPublisher) SetAuth(u string, pa string) {
	p.user, p.passwd = u, pa
}

func (p *kafkaPublisher) Run() error {
	if p.Opened() {
		return fmt.Errorf("kafka already run")
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = 5 * time.Second

	switch p.partition {
	case Manual:
		config.Producer.Partitioner = sarama.NewManualPartitioner
	case Random:
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	default:
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	//config.Producer.Retry =
	config.Version = sarama.V0_11_0_2

	if len(p.user) > 0 && len(p.passwd) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = p.user
		config.Net.SASL.Password = p.passwd
	}

	pro, err := sarama.NewAsyncProducer(p.hosts, config)
	if err != nil {
		return fmt.Errorf("kafkaPublisher init %v", err)
	}

	go func() {
		for {
			select {
			case _ = <-pro.Successes():

			case fail := <-pro.Errors():
				metric.IncCountWithClear(metric.Kafka)
				if fail == nil {
					logx.T(nil, define.ArchError).Errorf("kafka err response fail is nil")
					continue
				}

				logx.T(nil, define.ArchError).Errorf("kafka err response fail err:%v", fail.Err)
				// 1. the network jitter, causing messages not to be sent at all
				// 2. the message itself is not compliant, causing the broker to reject it
				// in short, the Producer is responsible for handling the failed send rather than brokers
				// if the broker is down at this point, the ticker will check that.
				p.failFuncCallback(fail.Msg, fail.Err)
			case <-p.close:
				return
			}
		}
	}()
	p.producer = pro
	return nil
}

// info[0] is key
// info[1] is partitionWay
func (p *kafkaPublisher) Publish(topic string, in io.Reader, info ...interface{}) (err error) {
	if !p.Running() {
		return fmt.Errorf("kafka not init yet")
	}
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	if len(info) == 0 {
		return fmt.Errorf("info need at least key")
	}

	// send message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(body),
	}

	if len(info) == 1 {
		msg.Key = sarama.StringEncoder(cast.ToString(info[0]))
	}

	if len(info) == 2 {
		msg.Key = sarama.StringEncoder(cast.ToString(info[0]))
		msg.Partition = cast.ToInt32(cast.ToInt(info[1]))
	}

	p.producer.Input() <- msg
	return
}

func (p *kafkaPublisher) PublishWithRetry(topic string, body []byte, try int32) error {
	if try <= 0 {
		logx.L().Warnf("kafkaPublisher try excess")
		return nil
	}
	err := p.Publish(topic, bytes.NewBuffer(body))
	if err != nil {
		atomic.AddInt32(&try, -1)
		logx.L().Warnf("kafkaPublisher untilTryout err: %s, topic: %s, leftTime: %v", err.Error(), topic, try)
		return p.PublishWithRetry(topic, body, try)
	}
	return err

}
