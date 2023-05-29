// Author: huaxinrui@tal.com
// Time: 2022-10-21 17:25
// Git: huaxr

package dispatch

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/huaxr/framework/pkg/toolutil"

	"github.com/Shopify/sarama"

	"github.com/huaxr/framework/component/dispatch/consumer"
	"github.com/huaxr/framework/component/dispatch/publisher"

	"github.com/huaxr/framework/logx"
)

var kafkaHosts = strings.Split("10.90.73.26:9092,10.90.73.54:9092,10.90.73.56:9092", ",")

var nsqHosts = strings.Split("10.90.72.135:4150,10.90.72.136:4150,10.90.72.171:4150,10.90.72.172:4150", ",")
var lookups = strings.Split("10.90.72.58:4161,10.90.72.94:4161", ",")
var admin = []string{"10.90.72.135:4171"}
var secret = "4BFE467B-FCBA-4519-BAC8-E9A3C57EDEB6"
var nsqTopic = "test"

const kafkaTopic = "xesFlow_first_example"

func TestKafkaPublish(t *testing.T) {
	callbackErr := func(m *sarama.ProducerMessage, e error) {
		t.Log(m.Key, e.Error())
	}
	pub := publisher.NewKafkaPublisher(context.Background(), kafkaHosts, publisher.Manual, callbackErr)
	if err := pub.Run(); err != nil {
		panic(err)
	}
	for {
		in := bytes.NewBuffer([]byte(toolutil.GetRandomString(10)))
		err := pub.Publish(kafkaTopic, in, "xx", 1)
		if err != nil {
			logx.L().Error(err)
		}
		time.Sleep(1 * time.Second)
	}
}

// 必须定义指定 partition
func TestKafkaConsumer(t *testing.T) {
	sub := consumer.NewKafkaConsumer(context.Background(), kafkaTopic, kafkaHosts, 1, true, func(msg *sarama.ConsumerMessage) {
		logx.L().Infof("%v %v %v", string(msg.Key), string(msg.Value), msg.Offset)
	})
	sub.SetAuth("", "")
	if err := sub.Run(); err != nil {
		panic(err)
	}
	sub.Consume()
}

func TestKafkaGroupConsumer(t *testing.T) {
	sub := consumer.NewKafkaGroupConsumer(context.Background(), []string{kafkaTopic}, kafkaHosts, "AAA", false, func(msg *sarama.ConsumerMessage) {
		logx.L().Infof("%v %v %v", string(msg.Key), string(msg.Value), msg.Offset)
	})
	if err := sub.Run(); err != nil {
		panic(err)
	}
	sub.Consume()
}

func TestNsqPublish(t *testing.T) {
	pub := publisher.NewNSQPublisher(context.Background(), nsqHosts, lookups, secret)
	pub.SetAuth("", "")
	if err := pub.Run(); err != nil {
		panic(err)
	}
	for {
		in := bytes.NewBuffer([]byte("HHHHHHHHHHHHHHHHHHHHHHHHHHHH"))
		err := pub.Publish(nsqTopic, in)
		if err != nil {
			logx.L().Error(err)
		}
		time.Sleep(2 * time.Second)
	}
}
