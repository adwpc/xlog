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
	DefaultTinyLogTimeFormat = "060102.150405.000"

	TraceLevel = "trace"
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

	//msg key+spliter+value
	OutputFormatNormalSpliter string
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

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

var TinyCoderConsoleConfig = Config{
	FileMode:       OutputToConsole,
	FormatMode:     OutputFormatNormal,
	CallerCodeLine: true,
	Level:          DebugLevel,
	Format:         DefaultTinyLogTimeFormat,
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func Init(c Config) {
	l := zerolog.GlobalLevel()
	switch c.Level {
	case TraceLevel:
		l = zerolog.TraceLevel
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
			return fmt.Sprintf("[%s:%d][%s] %v", fileName, line, funcName, i)
		}
		if c.CallerCodeLine && !c.CallerFuncName {
			return fmt.Sprintf("[%s:%d] %v", fileName, line, i)
		}
		if !c.CallerCodeLine && c.CallerFuncName {
			return fmt.Sprintf("[%s] %v", funcName, i)
		}
		return fmt.Sprintf("%v", i)
	}

	if c.FileMode == OutputToFile && c.FormatMode == OutputFormatJson {
		writer := initRollingWriter(c)
		logger = zerolog.New(writer).Level(l).With().Timestamp().Logger()
		if c.CallerCodeLine {
			logger = logger.With().Caller().Logger()
		}
		return
	}

	spliter := "=>"
	if c.OutputFormatNormalSpliter != "" {
		spliter = c.OutputFormatNormalSpliter
	}

	formatFieldName := func(i interface{}) string {
		return fmt.Sprintf("%s%s", i, spliter)
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
		output.FormatFieldName = formatFieldName

		logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
		return
	}

	if c.FileMode == OutputToConsole && c.FormatMode == OutputFormatNormal {
		formatFieldName := func(i interface{}) string {
			s := colorize(fmt.Sprintf("%v%v", i, spliter), colorCyan, false)
			return s
		}

		formatMessageFunc = func(i interface{}) string {
			caller, file, line, _ := runtime.Caller(8)
			fileName := filepath.Base(file)
			funcName := strings.TrimPrefix(filepath.Ext((runtime.FuncForPC(caller).Name())), ".")
			if c.CallerCodeLine && c.CallerFuncName {
				n := colorize(fmt.Sprintf("[%s:%d]", fileName, line), colorMagenta, false)
				l := colorize(fmt.Sprintf("[%s]", funcName), colorYellow, false)
				return fmt.Sprintf("%v%v %v", n, l, i)
			}
			if c.CallerCodeLine && !c.CallerFuncName {
				n := colorize(fmt.Sprintf("[%v:%d]", fileName, line), colorMagenta, false)
				return fmt.Sprintf("%v %v", n, i)
			}
			if !c.CallerCodeLine && c.CallerFuncName {
				l := colorize(fmt.Sprintf("[%v]", funcName), colorYellow, false)
				return fmt.Sprintf("%v %v", l, i)
			}
			return fmt.Sprintf("%v", i)
		}
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: c.Format}
		output.FormatTimestamp = func(i interface{}) string {
			return colorize(fmt.Sprintf("[%s]", i), colorDarkGray, false)
		}
		output.FormatLevel = func(i interface{}) string {
			level := i.(string)
			var l string
			switch level {
			case zerolog.LevelTraceValue:
				l = colorize("[T]", colorMagenta, false)
			case zerolog.LevelDebugValue:
				l = colorize("[D]", colorYellow, false)
			case zerolog.LevelInfoValue:
				l = colorize("[I]", colorGreen, false)
			case zerolog.LevelWarnValue:
				l = colorize("[W]", colorRed, false)
			case zerolog.LevelErrorValue:
				l = colorize(colorize("[E]", colorRed, false), colorBold, false)
			case zerolog.LevelFatalValue:
				l = colorize(colorize("[F]", colorRed, false), colorBold, false)
			case zerolog.LevelPanicValue:
				l = colorize(colorize("[P]", colorRed, false), colorBold, false)
			default:
				l = colorize(level, colorBold, false)
			}
			return l
		}
		output.FormatMessage = formatMessageFunc
		output.FormatFieldName = formatFieldName

		logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
	}

	if c.FileMode == OutputToConsole && c.FormatMode == OutputFormatJson {
		logger = zerolog.New(os.Stdout).Level(l).With().Timestamp().Logger()
		if c.CallerCodeLine {
			logger = logger.With().Caller().Logger()
		}
	}
}

func Tracef(format string, v ...interface{}) {
	logger.Trace().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
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
		if v == nil {
			// for reflect.TypeOf(nil) panic
			keysAndValues[i] = fmt.Sprintf("%+v", v)
			continue
		}
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			keysAndValues[i] = fmt.Sprintf("%+v", v)
		}
		if reflect.TypeOf(v).Kind() == reflect.UnsafePointer {
			keysAndValues[i] = fmt.Sprintf("%p", v)
		}
		switch v.(type) {
		case error:
			keysAndValues[i] = fmt.Sprintf("%+v", v.(error))
		}
	}
}

func Trace(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Trace().Fields(keysAndValues).Msg(msg)
}

func Debug(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Debug().Fields(keysAndValues).Msg(msg)
}

func Info(msg string, keysAndValues ...interface{}) {
	FormatStruct(keysAndValues...)
	logger.Info().Fields(keysAndValues).Msg(msg)
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
