## 特性 [English](https://github.com/adwpc/xlog/blob/master/README.md).

- [x] 使用简便
- [x] 基于 zerolog 的高性能
- [x] 支持多种打印字段
	- [x] 日志级别
	- [x] 调用者代码行（例如：example.go:32，你可以通过在 IDE 中按下 Alt 键并单击它来跳转到该行）
	- [x] 调用者函数名
	- [x] 时间格式
	- [x] 格式化消息（使用Infof/Debugf/Xxxxf 类似于 fmt.Printf()）
	- [x] 键值对消息（使用Info/Debug/Xxxx，推荐，无需编写 %v等）
	- [x] 支持所有类型 (int/string/float64/map/struct/interface)
	
- [x] 支持输出到控制台

	- [x] 支持输出普通格式
	
	```
	[2023-08-24 17:48:17.227] [I] [example.go:88][TestOutputToConsoleNormal] TestOutputToConsoleNormal Info float=1.23 int=123 map={"a":1,"b":2} string=abc struct="{a:1 b:{b:1} c:map[a:1 b:2]}"
	[2023-08-24 17:48:17.227] [I] [example.go:88][TestOutputToConsoleNormal] TestOutputToConsoleNormal Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}
	```
	
	
	- [x] 支持输出JSON格式
	
	```
	{"level":"info","time":"2023-08-24 18:12:04.099","caller":"example.go:32","message":"TestOutputToConsoleJson Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}"}
	{"level":"info","string":"abc","int":123,"float":1.23,"map":{"a":1,"b":2},"struct":"{a:1 b:{b:1} c:map[a:1 b:2]}","time":"2023-08-24 18:12:04.099","caller":"example.go:33","message":"TestOutputToConsoleJson Info"}
	```

- [x] 支持输出到文件
	- [x] 支持文件切割滚动
	
	```
	TestOutputToFileJson.log
	TestOutputToFileJson.log.2023-08-24 17:10:44
	TestOutputToFileJson.log.2023-08-24 17:30:03
	```
	
	- [x] 支持输出普通格式
	
	```
	[2023-08-24 17:30:08.343] [I] [example.go:75][TestOutputToFileNormal] TestOutputToFileNormal Info float=1.23 int=123 map={"a":1,"b":2} string=abc struct="{a:1 b:{b:1} c:map[a:1 b:2]}"
	[2023-08-24 17:30:08.343] [I] [example.go:75][TestOutputToFileNormal] TestOutputToFileNormal Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}
	```
	
	
	- [x] 支持输出JSON格式
	
	```
	{"level":"info","string":"abc","int":123,"float":1.23,"map":{"a":1,"b":2},"struct":"{a:1 b:{b:1} c:map[a:1 b:2]}","time":"2023-08-24 18:13:10.455","caller":"example.go:33","message":"TestOutputToFileJson Info"}
	{"level":"info","time":"2023-08-24 18:13:10.455","caller":"example.go:32","message":"TestOutputToFileJson Infof: abc 123 1.23 map[a:1 b:2] {a:1 b:{b:1} c:map[a:1 b:2]}"}

	```

