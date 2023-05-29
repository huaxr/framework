// Author: huaxr
// Time:   2022/1/26 下午2:47
// Git:    huaxr

package publisher

import (
	"fmt"

	"github.com/huaxr/framework/logx"

	"strings"
	"testing"
	"time"

	"github.com/huaxr/framework/pkg/toolutil"
	"github.com/Shopify/sarama"
)

const (
	topic = "xesFlow_first_example"
)

var hosts = strings.Split("brokers", ",")

// var hosts = []string{"10.90.73.71:9092", "10.90.73.74:9092", "10.90.73.89:9092"}

func Producer() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	producer, err := sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		fmt.Printf("producer_test create nsqProducer error :%s\n", err.Error())
		return
	}

	defer producer.AsyncClose()

	for {
		// send message
		msg := &sarama.ProducerMessage{
			Topic: topic,
			// hash this key
			Key:       sarama.StringEncoder("go_test"),
			Value:     sarama.ByteEncoder(toolutil.GetRandomString(3)),
			Partition: 10,
		}

		// send to chain
		producer.Input() <- msg

		logx.L().Infof("send: %v, partition:%v", msg.Value, msg.Partition)
		select {
		case suc := <-producer.Successes():
			fmt.Printf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
		case fail := <-producer.Errors():
			fmt.Printf("err: %s\n", fail.Err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}

func Consumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	// consumer
	consumer, err := sarama.NewConsumer(hosts, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}

	defer consumer.Close()

	partitionconsumer, err := consumer.ConsumePartition(topic, 6, sarama.OffsetOldest)

	//message := consumer.Messages()
	if err != nil {
		fmt.Printf("try create partition_consumer error %s\n", err.Error())
		return
	}
	defer partitionconsumer.Close()

	for {
		select {
		case msg := <-partitionconsumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-partitionconsumer.Errors():
			fmt.Printf("err :%s\n", err.Error())
		}
	}

}

func Metadata() {
	fmt.Printf("metadata tmp\n")

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2

	client, err := sarama.NewClient(hosts, config)
	if err != nil {
		fmt.Printf("metadata_test try create client err :%s\n", err.Error())
		return
	}

	defer client.Close()

	// get topic set
	topics, err := client.Topics()
	if err != nil {
		fmt.Printf("try get topics err %s\n", err.Error())
		return
	}

	fmt.Printf("topics(%d):\n", len(topics))

	for _, topic := range topics {
		fmt.Println(topic)
	}

	// get broker set
	brokers := client.Brokers()
	fmt.Printf("broker set(%d):\n", len(brokers))
	for _, broker := range brokers {
		fmt.Printf("%s\n", broker.Addr())
	}

	res, err := client.Partitions(topic)
	fmt.Printf("Partitions set(%v), err:%v", res, err)

}

func TestMetadata(t *testing.T) {
	Metadata()
}

func TestProducer(t *testing.T) {
	Producer()
}

func TestConsumer(t *testing.T) {
	Consumer()
}
