package ginx

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/huaxr/framework/version"

	"github.com/huaxr/framework/pkg/pprofutil"

	"github.com/huaxr/framework/component/promethu"
	"github.com/huaxr/framework/component/tcm"
	"github.com/huaxr/framework/component/ticker"
	"github.com/huaxr/framework/ginx/middleware"
	"github.com/huaxr/framework/ginx/response"
	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/huaxr/framework/pkg/httputil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	_ "go.uber.org/automaxprocs"
)

type Gin struct {
	define.Opening
	ctx         context.Context
	engine      *gin.Engine
	once        *sync.Once
	openMonitor bool
	openTcm     bool
	openPprof   bool
}

func NewGinx(ctx context.Context) *Gin {
	var b = make(chan string, 100)
	e := new(Gin)
	e.engine = gin.Default()
	// the middleware should be in sequence
	e.engine.Use(
		middleware.TimeMetric(),
		middleware.Limiter(),
		middleware.MarkPassThrough(),
		middleware.WrapCache(),
		gin.CustomRecoveryWithWriter(&logger{
			b: b,
		}, func(c *gin.Context, err interface{}) {
			select {
			case s := <-b:
				logx.T(c, define.ArchFatal).Errorf("gin server panic:%v ========== method:%v  ========== url:%v",
					s, c.Request.Method, c.Request.RequestURI)
			default:
			}
			response.Error(c, fmt.Errorf("montage internal error:%v", err))
			c.Abort()
		}),
	)

	e.ctx = ctx
	e.once = &sync.Once{}
	e.openMonitor = true
	e.openTcm = true
	e.openPprof = true
	return e
}

func (r *Gin) preStop() {

}

func (r *Gin) tcm() {
	if r.openTcm {
		tcm.InitTcmInstance(r.ctx)
	}
}

func (r *Gin) pprof() {
	if r.openPprof {
		pprofutil.GinPprof(r.engine)
	}
}

func (r *Gin) monitor() {
	if r.openMonitor {
		ticker.RegisterTick(&promethu.Generator{
			Cli: httputil.NewHttpClient(3*time.Second, 0, false),
		})
	}
}

func (r *Gin) Run() error {
	if r.Opened() {
		return fmt.Errorf("gin service already Run, do not call Run multi time")
	}

	ctx, cancel := context.WithCancel(r.ctx)
	r.ctx = ctx
	r.tcm()
	r.monitor()
	r.pprof()
	//router.LoadHTMLGlob("static/templates/*")
	//e.Router.StaticFS("/static", http.Dir("static"))
	mode := confutil.GetDefaultConfig().Gin.Mode
	switch mode {
	case confutil.Debug, confutil.Tester, confutil.Release:
		gin.SetMode(mode.String())
	default:
		gin.SetMode(confutil.Release.String())
	}
	gin.DefaultWriter = ioutil.Discard

	r.once = new(sync.Once)
	port := confutil.GetDefaultConfig().Gin.Port
	if err := define.ValidatePort(port); err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)
	// when pod deploy on sub-process, SIGTERM from k8s will not received by this channel
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		r.preStop()
		time.Sleep(1 * time.Second)
	}()

	host, _ := os.Hostname()
	msg := fmt.Sprintf(`--------------------------------------
   *     ^__^                        |- Welcome %s To Use Montage Ginx Framework(%s)!!
  \|/    (oo)\_______                |- Listening on    :%d    
         (__)\       )\/\            |- Basic Psm        %s
       *     ||----w |         ~~~~  |- Use Tcm          %s
      \|/    ||     ||      * U      |- Use Monitor      %s
            ~--     --~    \|/       |- Use Pprof        %s
 ~~~~~ ~~~~~~ ~~~~~~~~~~ ~~~~        |- Warning! Do Not Use The Same PSM For Other Different Services
--------------------------------------
`, host, version.GetVersionStr(), port, confutil.GetDefaultConfig().PSM.String(), cast.ToString(r.openTcm), cast.ToString(r.openMonitor), cast.ToString(r.openPprof))
	fmt.Print(msg)

	err := r.engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		cancel()
		return err
	}
	return nil
}

func (r *Gin) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Gin) RegisterMiddleware(middles ...gin.HandlerFunc) {
	// by default gin.DefaultWriter = os.Stdout
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.once.Do(func() {
		for _, m := range middles {
			r.engine.Use(m)
		}
	})
}

func (r *Gin) RegisterRouter(routers ...Router) {
	for _, i := range routers {
		i.Router(r.engine)
	}
}

func (r *Gin) RegisterPanicHook(f func(string)) {
	setRecoverHook(f)
}
