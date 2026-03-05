package logx

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

const (
	Panic uint32 = iota
	Fatal
	Error
	Warn
	Info
	Debug
	Trace
)

type Config struct {
	FileName   string
	Level      uint32
	MaxSize    int
	MaxBackups int
}

type Option func(*Config)

func WithLogName(name string) Option {
	return func(c *Config) {
		c.FileName = name
	}
}

func WithLevel(level uint32) Option {
	return func(c *Config) {
		c.Level = level
	}
}

func WithMaxSize(maxsize int) Option {
	return func(c *Config) {
		c.MaxSize = maxsize
	}
}

func WithMaxBackups(maxbackups int) Option {
	return func(c *Config) {
		c.MaxBackups = maxbackups
	}
}

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.000")
	var file, function string
	var line int
	if entry.Caller != nil {
		file = entry.Caller.File //filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
		function = entry.Caller.Function
		if funcs := strings.SplitN(function, ".", 2); len(funcs) == 2 {
			function = funcs[1]
		}
	}
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	case logrus.InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	}
	msg := fmt.Sprintf("%s \x1b[%dm%s/%s:%d %s\x1b[0m -- %s\n", timestamp, levelColor, strings.ToUpper(entry.Level.String()[:1]), file, line, function, entry.Message)
	return []byte(msg), nil
}

func New(opts ...Option) *logrus.Logger {
	c := &Config{
		FileName:   "log",
		Level:      Info,
		MaxSize:    4,
		MaxBackups: 10,
	}
	for _, opt := range opts {
		opt(c)
	}

	writer := &lumberjack.Logger{
		Filename:   c.FileName,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		LocalTime:  true,
	}
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetOutput(writer)
	log.SetFormatter(new(LogFormatter))

	return log
}
