[中文](https://github.com/adwpc/xlog/blob/master/README_CN.md)
## Getting Started

```
package main

import (
	xlog "github.com/adwpc/xlog"
)

func main() {
	xlog.Init(xlog.TinyCoderConsoleConfig)
	defer xlog.Close()

	xlog.Infof("Call Infof: %s %d %v", "abc", 123, 1.23)
	xlog.Info("Call Info", "string", "abc", "int", 123, "float", 1.23)
}
```

## Feature 
- [x] Easy to use
- [x] High performance base on zerolog
- [x] Support printing field
	- [x] log level
	- [x] caller code line (eg:`example.go:32` you can go to this line by alt+click it in IDE)
	- [x] caller function name
	- [x] time format
	- [x] format message (use Infof/Debugf/Xxxxf like fmt.Printf())
	- [x] key-value message (use Info/Debug/Xxxx recommend, without writing %v)
	- [x] all type support(int/string/float64/map/struct/interface)
	
- [x] Support output to console
	- [x] pretty logging
            ![tiny](https://raw.githubusercontent.com/adwpc/xlog/master/snapshot/tiny_console.png)
	- [x] Support output format normal
	
	```
	[2023-08-24 17:48:17.227] [I] [example.go:88][TestOutputToConsoleNormal] TestOutputToConsoleNormal Info float=1.23 int=123 map={"a":1,"b":2} string=abc struct="{a:1 b:{b:1} c:map[a:1 b:2]}"
	[2023-08-24 17:48:17.227] [I] [example.go:88][TestOutputToConsoleNormal] TestOutputToConsoleNormal Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}
	```
	
	
	- [x] Support output format json
	
	```
	{"level":"info","time":"2023-08-24 18:12:04.099","caller":"example.go:32","message":"TestOutputToConsoleJson Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}"}
	{"level":"info","string":"abc","int":123,"float":1.23,"map":{"a":1,"b":2},"struct":"{a:1 b:{b:1} c:map[a:1 b:2]}","time":"2023-08-24 18:12:04.099","caller":"example.go:33","message":"TestOutputToConsoleJson Info"}
	```

- [x] Support output to file
	- [x] Support file rolling
	
	```
	TestOutputToFileJson.log
	TestOutputToFileJson.log.2023-08-24 17:10:44
	TestOutputToFileJson.log.2023-08-24 17:30:03
	```
	
	- [x] Support output format normal
	
	```
	[2023-08-24 17:30:08.343] [I] [example.go:75][TestOutputToFileNormal] TestOutputToFileNormal Info float=1.23 int=123 map={"a":1,"b":2} string=abc struct="{a:1 b:{b:1} c:map[a:1 b:2]}"
	[2023-08-24 17:30:08.343] [I] [example.go:75][TestOutputToFileNormal] TestOutputToFileNormal Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}
	```
	
	
	- [x] Support output format json
	
	```
	{"level":"info","string":"abc","int":123,"float":1.23,"map":{"a":1,"b":2},"struct":"{a:1 b:{b:1} c:map[a:1 b:2]}","time":"2023-08-24 18:13:10.455","caller":"example.go:33","message":"TestOutputToFileJson Info"}
	{"level":"info","time":"2023-08-24 18:13:10.455","caller":"example.go:32","message":"TestOutputToFileJson Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}"}

	```

