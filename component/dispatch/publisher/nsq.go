// Author: huaxr
// Time:   2021/6/9 下午3:50
// Git:    huaxr

package publisher

import "C"
import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/huaxr/framework/component/plugin/selector"

	"io/ioutil"
	"log"
	"unsafe"

	"github.com/huaxr/framework/logx"
	"github.com/nsqio/go-nsq"
)

var (
	nsqProducer *nsqPublisher
	nsqLogger   = log.New(ioutil.Discard, "", log.LstdFlags)
)

// nsqPublisher provides function to dispense message to specified topic
// it's only works on server dispatch engine.
type nsqPublisher struct {
	sync.RWMutex
	ctx context.Context
	// binding address with client
	availableAddress map[string]unsafe.Pointer
	lookups          []string
	hosts            []string
	secret           string
	selector         selector.Selector
}

func NewNSQPublisher(ctx context.Context, hosts []string, lookups []string, secret string) Publisher {
	nsqProducer = new(nsqPublisher)
	nsqProducer.ctx = ctx
	nsqProducer.lookups = lookups
	nsqProducer.hosts = hosts
	nsqProducer.secret = secret
	nsqProducer.availableAddress = make(map[string]unsafe.Pointer)
	nsqProducer.selector = selector.NewSelector(selector.RoundRobin, hosts)
	return nsqProducer
}

// DELETE http://10.90.72.135:4171/api/nodes/10.90.72.172%3A4151
// BODY: {"topic": "xesFlow_first_example"}
func (p *nsqPublisher) delete(broker string, topic string) {}

func (p *nsqPublisher) SetAuth(u string, pa string) {}

func (p *nsqPublisher) Close() error {
	p.Lock()
	defer p.Unlock()
	if len(p.availableAddress) == 0 {
		return fmt.Errorf("no avaliable address found")
	}
	for _, pr := range p.availableAddress {
		pr2 := (*nsq.Producer)(pr)
		if pr2.Ping() != nil {
			pr2.Stop()
		}
	}

	p.availableAddress = map[string]unsafe.Pointer{}
	return nil
}

func (p *nsqPublisher) Run() error {
	for _, addr := range p.hosts {
		config := nsq.NewConfig()
		config.AuthSecret = p.secret
		pro, err := nsq.NewProducer(addr, config)
		if err != nil {
			logx.L().Errorf("nsqPublisher err %v", err)
			return err
		}
		// pro.ping will has nsq log
		pro.SetLogger(nsqLogger, nsq.LogLevelInfo)
		if err = pro.Ping(); err != nil {
			logx.L().Warnf("addr:%v ping err: %v", pro.String(), err)
			return err
		}
		p.availableAddress[pro.String()] = unsafe.Pointer(pro)
		logx.L().Debugf("start(recover) %s nsq publisher", pro.String())
	}
	return nil
}

// Publish dispatch AllowedBrokers dispatch message with the specified brokers
// with randomise & roundRobin strategy. the brokers should reserve quota by the app.
// familiar with knit one row, purl one row.
func (p *nsqPublisher) Publish(topic string, in io.Reader, infos ...interface{}) (err error) {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		logx.L().Errorf("nsq publish err:%v", err)
		return err
	}
	if len(body) == 0 || len(body) > 1048576 {
		return fmt.Errorf("dispatch err size: %v", len(body))
	}
	var count int
RETRY:
	broker := p.selector.Select()
	pp, _ := p.availableAddress[broker]
	err = (*nsq.Producer)(pp).Publish(topic, body)
	if err != nil {
		count++
		if count >= len(p.hosts) {
			logx.L().Errorf("publisher dead for %d times", count)
			return fmt.Errorf("nsq publish err:%v", err)
		}
		goto RETRY
	}

	return
}
