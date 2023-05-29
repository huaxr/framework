package logx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"github.com/huaxr/framework/version"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/pkg/confutil"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger         zLogger
	defaultLog     *zap.SugaredLogger
	defaultArchLog *zap.SugaredLogger
	// six
	logWithoutCtx *zap.SugaredLogger
)

type zLogger struct {
	// zap with basic fields
	levelLogger *zap.SugaredLogger
	archLogger  *zap.SugaredLogger
	pool        *sync.Pool
}

func WithArch() *zap.SugaredLogger {
	return defaultArchLog
}

func Flush() {
	_ = logger.levelLogger.Sync()
	_ = logger.archLogger.Sync()
}

func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.FullPath() + "]")
}

func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func init() {
	var encoder zapcore.Encoder
	encodeConfig := zapcore.EncoderConfig{
		MessageKey: define.Message.String(),
		LevelKey:   define.Level.String(),
		TimeKey:    define.Date.String(),
		CallerKey:  define.File.String(),
		EncodeTime: cEncodeTime,
	}
	switch confutil.GetDefaultConfig().Log.GetEncoder() {
	case confutil.Json:
		encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encodeConfig)

	case confutil.Console:
		encodeConfig.EncodeCaller = zapcore.FullCallerEncoder
		encodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encodeConfig)
	}

	l := confutil.GetDefaultConfig().Log.GetLogLevel()
	if l > zapcore.InfoLevel {
		fmt.Println("now not support warn and error, use info instead")
		l = zapcore.InfoLevel
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l
	})

	logPath := confutil.GetDefaultConfig().Log.Infofile
	infoWriter := getWriter(logPath)

	// define metric defaultL here
	path, fileName := filepath.Split(logPath)
	p := strings.Split(fileName, ".")
	if len(p) != 2 || p[1] != "log" {
		panic("defaultL suffix file path err, format: info.log")
	}
	// todo: split logs from metric from info
	//if strings.HasPrefix(fileName, "metric") {
	//	panic(fmt.Sprintf("fileName:%v is not allowed", fileName))
	//}
	metricPath := fmt.Sprintf("%s%s", path, fmt.Sprintf("%s.metric.log", p[0]))
	metricWriter := getWriter(metricPath)

	coreArch := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(metricWriter), infoLevel),
	}
	archL := zap.New(zapcore.NewTee(coreArch...), zap.AddCaller())

	coreDefault := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	}
	if confutil.GetDefaultConfig().Log.GetConsole() {
		coreDefault = append(coreDefault, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel))
	}
	defaultL := zap.New(zapcore.NewTee(coreDefault...), zap.AddCaller())

	logger = zLogger{
		levelLogger: defaultL.Sugar(),
		archLogger:  archL.Sugar(),
		pool: &sync.Pool{
			New: func() interface{} {
				return defineFields{}
			},
		},
	}

	defaultLog = logger.levelLogger.
		With(zap.String(define.Host.String(), ip.GetIp())).
		With(zap.String(define.PSM.String(), string(confutil.GetDefaultConfig().PSM))).
		With(zap.Int(define.Version.String(), version.Version))

	defaultArchLog = logger.archLogger.
		With(zap.String(define.Host.String(), ip.GetIp())).
		With(zap.String(define.PSM.String(), string(confutil.GetDefaultConfig().PSM))).
		With(zap.Int(define.Version.String(), version.Version))

	logWithoutCtx = defaultLog.
		With(zap.Float64(define.Duration.String(), -1)).
		With(zap.String(define.TraceId.String(), "")).
		With(zap.String(define.StartTime.String(), "")).
		With(zap.String(define.CallFrom.String(), "")).
		With(zap.String(define.HandlerPath.String(), "")).
		With(zap.Int64(define.HandlerExecutePeriod.String(), -1))
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+"%Y%m%d"+".log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func genSug(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return logWithoutCtx
	}

	var sug = defaultLog
	ctxData := logger.getFields(ctx)
	for _, i := range ctxData {
		if i.Key == "" || i.Type == zapcore.UnknownType {
			continue
		}
		sug = sug.With(i)
	}
	return sug
}
