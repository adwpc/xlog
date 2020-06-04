package xlog

import (
	"github.com/arthurkiller/rollingwriter"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

const (
	DefaultLogTimeFormat     = "2006-01-02 15:04:05.999"
	DefaultLogFileTimeFormat = "060102150405"

	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"

	VolumeRolling = rollingwriter.VolumeRolling
	TimeRolling   = rollingwriter.TimeRolling
)

type Config struct {
	Level         string
	Format        string
	Path          string //log path
	TimeTagFormat string //log file rotate tag
	File          string //log file name
	Rolling       int    //rolling policy: VolumeRolling or TimeRolling
	// - 时间滚动: 配置策略如同 crontable, 例如,每天0:0切分, 则配置 0 0 0 * * *
	// - 大小滚动: 配置单个日志文件(未压缩)的滚动大小门限, 入1G, 500M
	RollingPattern string //rolling pattern
	Remain         int    //rotate file remain total
	Gzip           bool   //gzip
}

func Init(c Config) {
	l := zerolog.GlobalLevel()
	switch c.Level {
	case DebugLevel:
		l = zerolog.DebugLevel
	case InfoLevel:
		l = zerolog.InfoLevel
	case WarnLevel:
		l = zerolog.WarnLevel
	case ErrorLevel:
		l = zerolog.ErrorLevel
	}

	zerolog.TimeFieldFormat = c.Format
	rc := rollingwriter.NewDefaultConfig()
	rc.LogPath = c.Path
	rc.TimeTagFormat = c.TimeTagFormat
	rc.FileName = c.File
	rc.MaxRemain = c.Remain
	rc.RollingPolicy = c.Rolling
	rc.RollingVolumeSize = c.RollingPattern
	rc.Compress = c.Gzip
	if rc.TimeTagFormat == "" {
		rc.TimeTagFormat = DefaultLogFileTimeFormat
	}
	writer, err := rollingwriter.NewWriterFromConfig(&rc)
	if err != nil {
		panic(err)
	}
	output := zerolog.ConsoleWriter{Out: writer, NoColor: true, TimeFormat: c.Format}
	log = zerolog.New(output).Level(l).With().Timestamp().Logger()
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}
