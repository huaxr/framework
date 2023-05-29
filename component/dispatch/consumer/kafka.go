// Author: XinRui Hua
// Time:   2022/1/27 下午2:44
// Git:    huaxr

package consumer

import (
	"context"
	"fmt"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/logx"
	"github.com/Shopify/sarama"
)

type consumerType int
type offsetWay string

const (
	group consumerType = iota
	partition
)

const (
	Oldest offsetWay = "oldest"
	Newest offsetWay = "newest"
)

type kafkaConsumer struct {
	ctx          context.Context
	user, passwd string
	topics       []string
	hosts        []string
	consumerType consumerType
	offsetWay    offsetWay

	partitionConsumer sarama.PartitionConsumer
	groupConsumer     sarama.ConsumerGroup

	handler func(*sarama.ConsumerMessage)
	isSync  bool

	// partition consumer
	partitionNumber int32
	// groupName consumer
	groupName string
}

// Setup when using groupName consumer
func (kafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (kafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h kafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// NewKafkaConsumer Specifying partition Consumption
// partitionNumber Specifies the partition to be bound
// isSync uses multiple goroutine consumption
// f is the handler function, wrapped by the user
func NewKafkaConsumer(ctx context.Context, topic string, hosts []string, partitionNumber int32, isSync bool, f func(*sarama.ConsumerMessage)) Consumer {
	return &kafkaConsumer{
		ctx:             ctx,
		topics:          []string{topic},
		hosts:           hosts,
		partitionNumber: partitionNumber,
		handler:         f,
		isSync:          isSync,
		consumerType:    partition,
	}
}

// NewKafkaGroupConsumer 组消费
func NewKafkaGroupConsumer(ctx context.Context, topics []string, hosts []string, groupName string, isSync bool, f func(*sarama.ConsumerMessage)) Consumer {
	return &kafkaConsumer{
		ctx:          ctx,
		topics:       topics,
		hosts:        hosts,
		handler:      f,
		isSync:       isSync,
		groupName:    groupName,
		consumerType: group,
	}
}

func (h *kafkaConsumer) SetAuth(u string, p string) {
	h.user = u
	h.passwd = p
}

func (h *kafkaConsumer) Run() error {
	var (
		err               error
		consumer          sarama.Consumer
		partitionConsumer sarama.PartitionConsumer
		groupConsumer     sarama.ConsumerGroup
	)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //
	config.Consumer.Offsets.Retry.Max = 3
	config.Version = sarama.V0_11_0_2

	if len(h.user) > 0 && len(h.passwd) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = h.user
		config.Net.SASL.Password = h.passwd
	}
	switch h.consumerType {
	case group:
		// using groupName
		// 开启自动提交，需要手动调用MarkMessage才有效
		config.Consumer.Offsets.AutoCommit.Enable = true
		groupConsumer, err = sarama.NewConsumerGroup(h.hosts, h.groupName, config)
		if err != nil {
			return err
		}
		h.groupConsumer = groupConsumer
		return nil
	case partition:
		// consumer
		consumer, err = sarama.NewConsumer(h.hosts, config)
		if err != nil {
			return err
		}
		partitionConsumer, err = consumer.ConsumePartition(h.topics[0], h.partitionNumber, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		h.partitionConsumer = partitionConsumer
		return nil
	default:
		return fmt.Errorf("not implement yet")
	}
}

func (h *kafkaConsumer) Consume() {
	switch h.consumerType {
	case group:
		_ = h.groupConsumer.Consume(h.ctx, h.topics, h)

	case partition:
		for {
			select {
			case msg := <-h.partitionConsumer.Messages():
				if h.isSync {
					go h.handler(msg)
					continue
				}
				h.handler(msg)

			case err := <-h.partitionConsumer.Errors():
				logx.L().Errorf("err :%s\n", err.Error())
				metric.IncCountWithClear(metric.Kafka)
			}
		}
	default:
		panic("not implement yet")
	}
}
