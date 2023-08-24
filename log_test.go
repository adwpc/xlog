package xlog

import (
	"testing"
	"time"
)

func Printing(tag string) {
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
			Infof(tag+" Infof: %s %d %v %v %+v", "abc", 123, 1.23, m, st)
			Info(tag+" Info", "string", "abc", "int", 123, "float", 1.23, "map", m, "struct", st)
		}
	}
}
func TestOutputToFileJson(t *testing.T) {
	c := Config{
		FileMode:           OutputToFile,
		FormatMode:         OutputFormatJson,
		Level:              DebugLevel,
		Format:             DefaultLogTimeFormat,
		FilePath:           "./logs",
		FileTimeTagFormat:  DefaultLogFileTimeFormat,
		FileName:           "TestOutputToFileJson",
		FileRollingPolicy:  VolumeRolling,
		FileRollingPattern: "10M",
		FileRemain:         3,
		FileUseGzip:        false,
	}

	Init(c)
	defer Close()
	Printing("TestOutputToFileJson")
}

func TestOutputToFileNormal(t *testing.T) {
	c := Config{
		FileMode:           OutputToFile,
		FormatMode:         OutputFormatNormal,
		Level:              DebugLevel,
		Format:             DefaultLogTimeFormat,
		FilePath:           "./logs",
		FileTimeTagFormat:  DefaultLogFileTimeFormat,
		FileName:           "TestOutputToFileNormal",
		FileRollingPolicy:  VolumeRolling,
		FileRollingPattern: "10M",
		FileRemain:         3,
		FileUseGzip:        false,
	}

	Init(c)
	defer Close()
	Printing("TestOutputToFileNormal")
}

func TestOutputToConsoleNormal(t *testing.T) {
	c := Config{
		FileMode:   OutputToConsole,
		FormatMode: OutputFormatNormal,
		Level:      DebugLevel,
		Format:     DefaultLogTimeFormat,
	}

	Init(c)
	defer Close()
	Printing("TestOutputToConsoleNormal")
}

func TestOutputToConsoleJson(t *testing.T) {
	c := Config{
		FileMode:   OutputToConsole,
		FormatMode: OutputFormatJson,
		Level:      DebugLevel,
		Format:     DefaultLogTimeFormat,
	}

	Init(c)
	defer Close()
	Printing("TestOutputToConsoleJson")

}
