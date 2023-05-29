// Author: XinRui Hua
// Time:   2022/3/22 下午3:38
// Git:    huaxr

package client

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/framework/pkg/toolutil"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"github.com/huaxr/framework/component/plugin/selector"
	"github.com/huaxr/framework/component/ticker"
	"github.com/huaxr/framework/internal/consensus"
	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/internal/metric"
	"github.com/huaxr/framework/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type conn struct {
	conn *grpc.ClientConn
	addr string
}

type service struct {
	define.Opening
	ctx context.Context
	// signal to tell weather services poll contaminate
	contaminate chan struct{}
	// set config to define such as retry dialOptions.
	dialConfig string
	// basic psm
	srvPath   string
	lock      sync.Mutex
	srv       map[string]conn
	ips       []string
	selector  selector.Selector
	consensus *consensus.Consensus
}

// when many pod restart or expansion, to shows
// some pod may lose some available pods ip, cause when Get from etcd
// and the starting watch goroutine got some lapse behind, we should
// keep watching it to avoiding some diff.
func (s *service) initSrv() {
	ips := s.consensus.Get(s.srvPath)
	if len(ips) == 0 {
		metric.Metric(define.HostEmpty, fmt.Sprintf("host empty for:%s", s.srvPath))
		return
	}

	for _, i := range ips {
		if err := s.register(i); err != nil {
			continue
		}
	}
}

// normally there would not any contaminate.
// if it really happens, auto repair it.
//
func (s *service) forceContaminate() {
	for {
		<-s.contaminate
		// pod restart cause lase
		metric.Metric(define.ServiceContaminate, fmt.Sprintf("ServiceContaminate, current srv size:%v, all:%v", len(s.srv), s))
		s.srv = make(map[string]conn)
		s.selector = selector.NewSelector(selector.RoundRobin, []string{})
		s.initSrv()
	}
}

func (s *service) checkDiff() {
	ips := s.consensus.Get(s.srvPath)
	if len(ips) == 0 {
		return
	}

	if len(ips) == len(s.srv) {
		// there are possible contaminate here, let getConn fix that lazy.
		// available above p999
		return
	}

	if len(ips) < len(s.srv) {
		metric.Metric(define.ServiceDiff, fmt.Sprintf("ServiceDiff, len(ips) < len(s.srv) current etcd keys:%v, current srv:%v", len(ips), len(s.srv)))
		select {
		case s.contaminate <- struct{}{}:
		default:
		}
		return
	}

	if len(ips) > len(s.srv) {
		// when restart the cluster, len(ips)=30 and s.srv size is 34. that's queer!
		// that's why we should preStop our service.
		metric.Metric(define.ServiceDiff, fmt.Sprintf("ServiceDiff, len(ips) > len(s.srv) current etcd keys:%v, current srv:%v", len(ips), len(s.srv)))
		for _, i := range ips {
			if err := s.register(i); err != nil {
				break
			}
		}
	}
}

func (s *service) TickHeartbeat() *ticker.T {
	// service available 100% guarantee thread.
	// like some etcd cluster panic cases.
	rand.Seed(time.Now().UnixNano())
	du := time.Duration(rand.Intn(10) + 50)
	return ticker.NewT(s.checkDiff, fmt.Sprintf("rpc:%s", s.srvPath), time.NewTicker(du*time.Second))
}

// service with naming index.
func newService(prefix string, cli *consensus.Consensus) *service {
	return &service{
		// the channel size can keep the concurrency when calling
		// the handle to deal with contaminate
		// set 0 will not work when default case set
		// 100 will not limit the execution, we should control etcd query qps.
		contaminate: make(chan struct{}, 1),
		srvPath:     prefix,
		lock:        sync.Mutex{},
		srv:         make(map[string]conn),
		selector:    selector.NewSelector(selector.RoundRobin, []string{}),
		ctx:         context.Background(),
		consensus:   cli,
	}
}

func (s *service) String() string {
	var addrs = make([]string, 0)
	for _, v := range s.srv {
		addrs = append(addrs, v.addr)
	}
	return "[" + strings.Join(addrs, ", ") + "]"
}

func (s *service) flush() {
	var alivePods = make([]string, 0, len(s.srv))
	for k := range s.srv {
		alivePods = append(alivePods, k)
	}
	ip.Sort(alivePods)
	s.ips = alivePods
	s.selector = selector.NewSelector(selector.RoundRobin, alivePods)
}

func (s *service) connect(addr string) (con *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	if len(s.dialConfig) > 0 {
		con, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(s.dialConfig))
	} else {
		con, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	}
	if err != nil {
		// when k8s services restart
		return nil, fmt.Errorf("dial addr %v err %v, psm:%v", addr, err, s.srvPath)
	}
	return
}

func (s *service) register(addr string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	con, err := s.connect(addr)
	if err != nil {
		logx.T(nil, define.ArchError).Infof("register addr %v dial err, err:%v, path:%v", addr, err, s.srvPath)
		return err
	}

	// when rpc server exit with panic, the pod will recreate the pid-1 to start deployment process,;;;;;;;;;;;;;;ll (which is not grace
	// stop) that intermittent restart lapse will cause serve watch an put before lease expire.
	if _, ok := s.srv[addr]; ok {
		logx.T(nil, define.ArchError).Infof("register addr %v already exist, err:%v, path:%v", addr, err, s.srvPath)
		return nil
	}

	s.srv[addr] = conn{
		conn: con,
		addr: addr,
	}

	s.flush()
	return nil
}

func (s *service) unregister(addr string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.srv[addr]; !ok {
		logx.T(nil, define.ArchError).Infof("unregister addr %v not exist, path:%v", addr, s.srvPath)
		return
	}
	// tcp handshakes on waves process:
	// syn_sent(c) --[syn=1,seq=x]--> syn_rcvd(s) --[ack=x+1,seq=y]-> established(c) --[ack=y+1]--> established(s)
	// after the connection is deleted from the map, syn_sent is sent again,
	// but the pod has been destroyed and syn_rcvd is not returned.
	// in this case, the tcp connection keeps holding the handle (netstat)
	// so we should close it on hand.
	err := s.srv[addr].conn.Close()
	if err != nil {
		logx.T(nil, define.ArchError).Infof("unregister addr %v  err:%v, path:%v", addr, err, s.srvPath)
	}
	delete(s.srv, addr)
	s.flush()
}

// if the method binding with service rather than *service.
// it will be shadowed cause. (*s) option is required here
// cause unregister changed the struct.
func (s *service) selectConn(addr string) (*grpc.ClientConn, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	c, ok := s.srv[addr]
	if !ok {
		return nil, fmt.Errorf("host not exist addr:%v", addr)
	}
	if c.conn.GetState() > connectivity.Ready {
		// maybe user close it by self
		con, err := s.connect(addr)
		if err != nil {
			metric.Metric(define.GrpcClosed)
			// do not unregister cause lock has been used
			delete(s.srv, addr)
			s.flush()
			return nil, err
		}
		s.srv[addr] = conn{
			conn: con,
			addr: addr,
		}
		return con, nil
	}
	return c.conn, nil
}

func (s *service) try() error {
	if !s.Running() {
		return fmt.Errorf("service not init yet")
	}
	if len(s.ips) == 0 {
		return fmt.Errorf("no avaliable host for:%v currently", s.srvPath)
	}
	return nil
}

func (s *service) GetConnByIp(ip string) (*grpc.ClientConn, error) {
	if err := s.try(); err != nil {
		return nil, err
	}
	return s.selectConn(ip)
}

// GetConn returns a persistent connection,
// You Must Not call conn.close when you use it to establish your data transform.
// otherwise, the connection is dead while you calling rpc method

// Must not keep this returning connection as global variable
// it actually ephemeral one through keepalive as a long connection.

// connection reuse, concurrent security, please consult test
func (s *service) GetConn() (*grpc.ClientConn, error) {
	if err := s.try(); err != nil {
		return nil, err
	}
	addr := s.selector.Select()
	conn, err := s.selectConn(addr)
	if err != nil {
		return s.GetConn()
	}
	return conn, nil
}

// if you want close your client when you fetch on established socket(defer con.close e.g.)
// using GetDialConn to selectConn a new one instead.
func (s *service) GetNewDialConn() (*grpc.ClientConn, error) {
	if err := s.try(); err != nil {
		return nil, err
	}
	addr := s.selector.Select()
	con, err := s.connect(addr)
	if err != nil {
		metric.Metric(define.GrpcClosed)
		s.unregister(addr)
		return nil, err
	}
	return con, nil
}

func (s *service) GetNewDialConnByIp(ip string) (*grpc.ClientConn, error) {
	if err := s.try(); err != nil {
		return nil, err
	}
	con, err := s.connect(ip)
	if err != nil {
		metric.Metric(define.GrpcClosed)
		s.unregister(ip)
		return nil, err
	}
	return con, nil
}

func (s *service) Run() error {
	if s.Opened() {
		return fmt.Errorf("service found client for %v already Run, do not call Run multi time", s.srvPath)
	}
	serviceLock.Lock()
	defer serviceLock.Unlock()
	servicePool[s.srvPath] = s
	go s.forceContaminate()
	go watcher(s, s.consensus.WatchKey(s.srvPath))
	s.initSrv()
	go ticker.RegisterTick(s)
	return nil
}

// grpc internal retry policy:
// 1. grpc.WithDefaultServiceConfig(retryPolicy)
// 2. export GRPC_GO_RETRY=on
func (s *service) SetRetryConfig(ss string) {
	if on := os.Getenv("GRPC_GO_RETRY"); on != "on" {
		logx.L().Warnf("SetRetryConfig would not work. please `export GRPC_GO_RETRY=on` first")
	}

	if len(ss) == 0 {
		return
	}
	s.dialConfig = ss
}

func (s *service) GetIpByModKey(key string) (string, error) {
	if err := s.try(); err != nil {
		return "", err
	}
	addr := s.ips[toolutil.CRC([]byte(key))%len(s.ips)]
	return addr, nil
}
