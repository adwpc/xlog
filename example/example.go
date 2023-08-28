package main

import (
	"time"
	"errors"

	xlog "github.com/adwpc/xlog"
)

func main() {
	TestOutputToFileJson()
	TestOutputToFileNormal()
	TestOutputToConsoleNormal()
	TestOutputToConsoleJson()
}

func TestPrinting(tag string) {
	ticker := time.NewTicker(time.Second * 3)
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
			xlog.Infof(tag+" Infof: %s %d %v %v %+v", "abc", 123, 1.23, m, st)
			xlog.Info(tag+" Info", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st)
			xlog.Error("msg:", errors.New("err"), "aa", errors.New("errr1"))
		}
	}
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
	TestPrinting("TestOutputToFileJson")
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
	TestPrinting("TestOutputToFileNormal")
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
	TestPrinting("TestOutputToConsoleNormal")
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
	TestPrinting("TestOutputToConsoleJson")
}
