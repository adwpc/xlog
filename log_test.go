package xlog

import "testing"

func TestFunction(t *testing.T) {
	c := Config{
		Level:          InfoLevel,
		Format:         DefaultLogTimeFormat,
		Path:           "./logs",
		TimeTagFormat:  DefaultLogFileTimeFormat,
		File:           "test",
		Rolling:        VolumeRolling,
		RollingPattern: "10M",
		Remain:         3,
		Gzip:           false,
	}

	Init(c)
	for {
		Infof("hello world!")
	}
}
