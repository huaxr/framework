// Author: XinRui Hua
// Time:   2022/3/18 下午4:43
// Git:    huaxr

package grpcx

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"github.com/huaxr/framework/version"

	"github.com/huaxr/framework/pkg/pprofutil"

	"github.com/huaxr/framework/component/plugin/apicache"
	"github.com/huaxr/framework/component/plugin/circuit"
	"github.com/huaxr/framework/component/promethu"
	"github.com/huaxr/framework/component/tcm"
	"github.com/huaxr/framework/component/ticker"
	"github.com/huaxr/framework/internal/consensus"
	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/internal/metric"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/huaxr/framework/pkg/httputil"
	"github.com/huaxr/framework/pkg/toolutil"
	"github.com/spf13/cast"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/util/wait"
)

var hook InterfereHook
var useCache = true

type Grpcx struct {
	define.Opening
	Srv         *grpc.Server
	ctx         context.Context
	servicePath string
	port        int
	openFound   bool
	openTcm     bool
	openMonitor bool
	openPprof   bool
}

// reload
func ctx2ctx(ctx context.Context) context.Context {
	var (
		ma metadata.MD
		ok bool
	)

	ma, ok = metadata.FromIncomingContext(ctx)
	// ma default has :authority,content-type,user-agent
	if !ok {
		logx.T(ctx, define.ArchError).Infof("ctx2ctx failed")
		return context.Background()
	}

	// we can't judge here whether the rpc request is from gin or from kafka-consumer
	// if normal gin->rpc forget using transport package to transforming ctx before calling grpcx,
	// we should ignore that and keep the trace linked here by redefine the pass through fields
	if _, ok := ma[define.TraceId.String()]; !ok {
		ctx = context.WithValue(ctx, define.TraceId.String(), define.Uid())
		ctx = context.WithValue(ctx, define.StartTime.String(), define.Nano())
		// e.g.  kafka consumer call rpc
		ctx = context.WithValue(ctx, define.CallFrom.String(), define.Unknown())

		defaultMd := metadata.MD{}
		for _, i := range define.CtxPassThrough {
			if v := ctx.Value(i.String()); v != nil {
				defaultMd.Set(i.String(), cast.ToString(v))
			}
		}
		// we must set mdOutgoingKey{} in order to keep rpc->rpc->rpc flow trace
		ctx = metadata.NewOutgoingContext(ctx, defaultMd)
		return ctx
	}

	for k, v := range ma {
		if len(v) > 0 {
			ctx = context.WithValue(ctx, k, v[0])
		}
	}
	return ctx
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	path := info.FullMethod // 	/proto.TestRpc/GetCache
	// time-consuming calc should define here instead of f1
	var t1 = time.Now()
	defer func() {
		var t2 = time.Now()
		metric.MetRpc(ctx, path, t2.Sub(t1).Milliseconds())
	}()

	// api-cache
	if useCache {
		v, ok := apicache.GetApiCache().GetCache().KVGet(path)
		if ok {
			return v, nil
		}
	}

	f1 := func() error {

		// handler panic should be caught here rather than
		// outside this enclosure method.
		defer func() {
			if e := recover(); e != nil {
				stack := toolutil.GetStack()
				hook.Panic(stack)
				logx.T(ctx, define.ArchFatal).Errorf("rpc handler panic:%v, stack:%v, rpc:%v", e, stack, info.FullMethod)
				err = fmt.Errorf("backend rpc panic: %v", e)
			}

			// todo: add global timeout control
			//select {
			//case <-ctx.Done():
			//
			//default:
			//
			//}
		}()

		// time duration should not bury here
		ctx = ctx2ctx(ctx)
		ctx = hook.Before(ctx)
		resp, err = handler(ctx, req)
		hook.After(ctx)

		// api-cache
		if vv, ok := apicache.GetApiCache().GetExpSec(path); ok && useCache && err != nil {
			apicache.GetApiCache().GetCache().KVSet(path, resp, time.Duration(vv)*time.Second)
		}
		return nil
	}

	f2 := func(e error) error {
		metric.MetRpcFusing(ctx, path)
		return e
	}

	// to err must be set here otherwise the client will never get the information
	errX := circuit.GetBreakerSet().Monitor(path, f1, f2)
	if errX != nil {
		err = errX
	}
	return
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer func() {
		if err := recover(); err != nil {
			logx.T(nil, define.ArchFatal).Errorf("StreamServerInterceptor rpc handler panic:%v", err)
		}
	}()
	err := handler(srv, ss)
	return err
}

func NewGrpcx(ctx context.Context) *Grpcx {
	server := grpc.NewServer(
		grpc.StreamInterceptor(StreamServerInterceptor),
		grpc.UnaryInterceptor(UnaryServerInterceptor),
		// connection reuse: logicConn = conn * MaxConcurrentStreams
		grpc.MaxConcurrentStreams(1<<14+1<<8),
	)
	port := confutil.GetDefaultConfig().Grpc.Port
	if err := define.ValidatePort(port); err != nil {
		panic(err)
	}
	srv := fmt.Sprintf("/%s", confutil.GetDefaultConfig().PSM)
	return &Grpcx{
		ctx:         ctx,
		Srv:         server,
		servicePath: srv,
		port:        port,
		openFound:   true,
		openTcm:     true,
		openMonitor: true,
		openPprof:   true,
	}
}

func (g *Grpcx) RegisterHook(h InterfereHook) {
	hook = h
}

func (g *Grpcx) founder() {
	if g.openFound {
		addr := ip.GetIp()
		key := fmt.Sprintf("%s/%s:%d", g.servicePath, addr, g.port)
		alive := consensus.KeepAlive{
			Key: key,
			Ttl: 10,
		}
		go wait.UntilWithContext(g.ctx, alive.Alive, time.Second*3)
		g.listenKill(addr)
	}
}

func (g *Grpcx) tcm() {
	if g.openTcm {
		tcm.InitTcmInstance(g.ctx)
	}
}

func (g *Grpcx) monitor() {
	if g.openMonitor {
		ticker.RegisterTick(&promethu.Generator{
			Cli: httputil.NewHttpClient(3*time.Second, 0, false),
		})
	}
}

func (g *Grpcx) pprof() {
	if g.openPprof {
		pprofutil.Pprof()
	}
}

func (g *Grpcx) Run() error {
	if g.Opened() {
		return fmt.Errorf("grpc service %v already Run, do not call Run multi time", g.servicePath)
	}

	if hook == nil {
		hook = defaultHook{}
	}
	g.tcm()
	g.founder()
	g.monitor()
	g.pprof()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		panic(err)
	}
	host, _ := os.Hostname()
	msg := fmt.Sprintf(`--------------------------------- 
   *     ^__^			|- Welcome %s To Use Montage Grpcx Framework(%s)!!
  \|/    (oo)\_______		|- Listening     :%d
         (__)\       )\/\	|- Basic Psm      %s
       *     ||----w |        	|- Use Tcm        %s
      \|/    ||     ||      *	|- Use Founder    %s
            ~--     --~    \|/	|- Use Monitor    %s
 ~~~~~ ~~~~~~ ~~~~~~~~~~ ~~~~	|- Use Pprof      %s
 ~~~~~ ~~~~~~ ~~~~~~~~~~ ~~~~   |- Warning! Do Not Use The Same PSM For Other Different Services
---------------------------------	
`, host, version.GetVersionStr(), confutil.GetDefaultConfig().Grpc.Port, fmt.Sprintf("%s", confutil.GetDefaultConfig().PSM), cast.ToString(g.openTcm),
		cast.ToString(g.openFound), cast.ToString(g.openMonitor), cast.ToString(g.openPprof))
	fmt.Print(msg)
	err = g.Srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
