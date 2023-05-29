// Author: XinRui Hua
// Time:   2022/4/5 上午11:05
// Git:    huaxr

package confutil

import "go.uber.org/zap/zapcore"

type EncodeType string

const (
	Console EncodeType = "console"
	Json    EncodeType = "json"
)

var defaultLog = &Log{
	Level:       "info",
	Encoder:     "console",
	Console:     true,
	Disabletags: true,
	Infofile:    "/tmp/info.log",
}

type Log struct {
	Level       string     `yaml:"level"`
	Encoder     EncodeType `yaml:"encoder"`
	Console     bool       `yaml:"console"`
	Disabletags bool       `yaml:"disabletags"`
	Infofile    string     `yaml:"infofile"`
}

func (l *Log) GetLogLevel() zapcore.Level {
	switch l.Level {
	default:
		return zapcore.DebugLevel
	case "info", "warn", "error":
		return zapcore.InfoLevel
	}
}

func (l *Log) GetEncoder() EncodeType {
	return l.Encoder
}

func (l *Log) GetConsole() bool {
	return l.Console
}
