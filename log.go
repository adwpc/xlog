package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/arthurkiller/rollingwriter"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

const (
	DefaultLogTimeFormat     = "2006-01-02 15:04:05.000"
	DefaultLogFileTimeFormat = "2006-01-02 15:04:05"

	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"

	VolumeRolling = rollingwriter.VolumeRolling
	TimeRolling   = rollingwriter.TimeRolling

	OutputToConsole = OutputFileMode(0)
	OutputToFile    = OutputFileMode(1)

	OutputFormatNormal = OutputFormatMode(0)
	OutputFormatJson   = OutputFormatMode(1)
)

var (
	writer rollingwriter.RollingWriter
)

type OutputFileMode int
type OutputFormatMode int

type Config struct {
	// log mode
	FileMode       OutputFileMode
	FormatMode     OutputFormatMode
	CallerCodeLine bool
	CallerFuncName bool

	// zerolog
	Level  string
	Format string

	// rollingwriter
	FilePath          string //log path
	FileTimeTagFormat string //log file rotate tag
	FileName          string //log file name
	FileRollingPolicy int    //rolling policy: VolumeRolling or TimeRolling
	// - 时间滚动: 配置策略如同 crontable, 例如,每天0:0切分, 则配置 0 0 0 * * *
	// - 大小滚动: 配置单个日志文件(未压缩)的滚动大小门限, 入1G, 500M
	FileRollingPattern string //rolling pattern
	FileRemain         int    //rotate file remain total
	FileUseGzip        bool   //gzip

}

func initRollingWriter(c Config) rollingwriter.RollingWriter {
	rc := rollingwriter.NewDefaultConfig()

	// maybe lost logs when async mode
	// rc.WriterMode = "async"
	rc.LogPath = c.FilePath
	rc.TimeTagFormat = c.FileTimeTagFormat
	rc.FileName = c.FileName
	rc.MaxRemain = c.FileRemain
	rc.RollingPolicy = c.FileRollingPolicy
	rc.RollingVolumeSize = c.FileRollingPattern
	rc.Compress = c.FileUseGzip
	if rc.TimeTagFormat == "" {
		rc.TimeTagFormat = DefaultLogFileTimeFormat
	}
	var err error
	writer, err = rollingwriter.NewWriterFromConfig(&rc)
	if err != nil {
		panic(err)
	}
	return writer
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
	zerolog.CallerSkipFrameCount = 3
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	formatMessageFunc := func(i interface{}) string {
		caller, file, line, _ := runtime.Caller(8)
		fileName := filepath.Base(file)
		funcName := strings.TrimPrefix(filepath.Ext((runtime.FuncForPC(caller).Name())), ".")
		if c.CallerCodeLine && c.CallerFuncName {
			return fmt.Sprintf("[%s:%d][%s] %s", fileName, line, funcName, i)
		}
		if c.CallerCodeLine && !c.CallerFuncName {
			return fmt.Sprintf("[%s:%d] %s", fileName, line, i)
		}
		if !c.CallerCodeLine && c.CallerFuncName {
			return fmt.Sprintf("[%s] %s", funcName, i)
		}
		return fmt.Sprintf("%s", i)
	}

	if c.FileMode == OutputToFile && c.FormatMode == OutputFormatJson {
		writer := initRollingWriter(c)
		logger = zerolog.New(writer).Level(l).With().Timestamp().Logger()
		if c.CallerCodeLine {
			logger = logger.With().Caller().Logger()
		}
		return
	}

	if c.FileMode == OutputToFile && c.FormatMode == OutputFormatNormal {
		writer := initRollingWriter(c)
		output := zerolog.ConsoleWriter{Out: writer, TimeFormat: c.Format, NoColor: true}
		output.FormatTimestamp = func(i interface{}) string {
			return fmt.Sprintf("[%s]", i)
		}
		output.FormatLevel = func(i interface{}) string {
			level := i.(string)
			return strings.ToUpper(fmt.Sprintf("[%s]", string(level[0])))
		}
		output.FormatMessage = formatMessageFunc

		logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
		return
	}

	if c.FileMode == OutputToConsole && c.FormatMode == OutputFormatNormal {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: c.Format}
		output.FormatTimestamp = func(i interface{}) string {
			return fmt.Sprintf("[%s]", i)
		}
		output.FormatLevel = func(i interface{}) string {
			level := i.(string)
			return strings.ToUpper(fmt.Sprintf("[%s]", string(level[0])))
		}
		output.FormatMessage = formatMessageFunc

		logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
	}

	if c.FileMode == OutputToConsole && c.FormatMode == OutputFormatJson {
		logger = zerolog.New(os.Stdout).Level(l).With().Timestamp().Logger()
		if c.CallerCodeLine {
			logger = logger.With().Caller().Logger()
		}
	}
}

func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Error().Msgf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panic().Msgf(format, v...)
}

func FormatStruct(keysAndValues ...interface{}) {
	for i, v := range keysAndValues {
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			keysAndValues[i] = fmt.Sprintf("%+v", v)
		}
		switch err.(type) {
		case error:
			keysAndValues[i] = fmt.Sprintf("%+v", v.(error))
		}
	}
}

func Info(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Info().Fields(keysAndValues).Msg(msg)
}

func Debug(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Debug().Fields(keysAndValues).Msg(msg)
}

func Warn(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Warn().Fields(keysAndValues).Msg(msg)
}

func Error(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Error().Fields(keysAndValues).Msg(msg)
}

func Panic(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Panic().Fields(keysAndValues).Msg(msg)
}

func Close() {
	if writer != nil {
		writer.Close()
	}
}
