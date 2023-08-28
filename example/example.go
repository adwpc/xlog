package main

import (
	"errors"
	"time"

	xlog "github.com/adwpc/xlog"
)

func main() {
	TestTinyCoderConsole()
	TestOutputToFileJson()
	TestOutputToFileNormal()
	TestOutputToConsoleNormal()
	TestOutputToConsoleJson()
}

func TestPrinting(tag string, second int64) {
	ticker := time.NewTicker(time.Second * time.Duration(second))
	for {
		select {
		case <-ticker.C:
			return
		default:
			type B struct {
				b int
			}
			st := struct {
				a int
				b B
				c map[string]interface{}
			}{a: 1, b: B{b: 1}, c: map[string]interface{}{"a": 1, "b": 2}}
			m := map[string]interface{}{"a": 1, "b": 2}
			err := errors.New("an error")
			xlog.Debugf(tag+" Debugf: %s %d %v %v %+v %v", "abc", 123, 1.23, m, st, err)
			xlog.Debug(tag+" Debug", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st, "err", err)
			xlog.Infof(tag+" Infof: %s %d %v %v %+v %v", "abc", 123, 1.23, m, st, err)
			xlog.Info(tag+" Info", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st, "err", err)
			xlog.Warnf(tag+" Warnf: %s %d %v %v %+v %v", "abc", 123, 1.23, m, st, err)
			xlog.Warn(tag+" Warn", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st, "err", err)
			xlog.Errorf(tag+" Errorf: %s %d %v %v %+v %v", "abc", 123, 1.23, m, st, err)
			xlog.Error(tag+" Error", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st, "err", err)
		}
	}
}

func TestTinyCoderConsole() {
	xlog.Init(xlog.TinyCoderConsoleConfig)
	defer xlog.Close()
	TestPrinting("TestTinyCoderConsole", 1)
}

func TestOutputToFileJson() {
	c := xlog.Config{
		FileMode:           xlog.OutputToFile,
		FormatMode:         xlog.OutputFormatJson,
		CallerCodeLine:     true,
		Level:              xlog.DebugLevel,
		Format:             xlog.DefaultLogTimeFormat,
		FilePath:           "./logs",
		FileTimeTagFormat:  xlog.DefaultLogFileTimeFormat,
		FileName:           "TestOutputToFileJson",
		FileRollingPolicy:  xlog.VolumeRolling,
		FileRollingPattern: "10M",
		FileRemain:         3,
		FileUseGzip:        false,
	}

	xlog.Init(c)
	defer xlog.Close()
	TestPrinting("TestOutputToFileJson", 5)
}

func TestOutputToFileNormal() {
	c := xlog.Config{
		FileMode:           xlog.OutputToFile,
		FormatMode:         xlog.OutputFormatNormal,
		CallerCodeLine:     true,
		CallerFuncName:     true,
		Level:              xlog.DebugLevel,
		Format:             xlog.DefaultLogTimeFormat,
		FilePath:           "./logs",
		FileTimeTagFormat:  xlog.DefaultLogFileTimeFormat,
		FileName:           "TestOutputToFileNormal",
		FileRollingPolicy:  xlog.VolumeRolling,
		FileRollingPattern: "10M",
		FileRemain:         3,
		FileUseGzip:        false,
	}

	xlog.Init(c)
	defer xlog.Close()
	TestPrinting("TestOutputToFileNormal", 5)
}

func TestOutputToConsoleNormal() {
	c := xlog.Config{
		FileMode:       xlog.OutputToConsole,
		FormatMode:     xlog.OutputFormatNormal,
		CallerCodeLine: true,
		CallerFuncName: true,
		Level:          xlog.DebugLevel,
		Format:         xlog.DefaultLogTimeFormat,
	}

	xlog.Init(c)
	defer xlog.Close()
	TestPrinting("TestOutputToConsoleNormal", 1)
}

func TestOutputToConsoleJson() {
	c := xlog.Config{
		FileMode:       xlog.OutputToConsole,
		FormatMode:     xlog.OutputFormatJson,
		CallerCodeLine: true,
		Level:          xlog.DebugLevel,
		Format:         xlog.DefaultLogTimeFormat,
	}

	xlog.Init(c)
	defer xlog.Close()
	TestPrinting("TestOutputToConsoleJson", 1)
}
