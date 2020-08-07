package logger


//ILogging 基础日志记录接口
type ILogging interface {
	Printf(format string, content ...interface{})
	Print(content ...interface{})
	Println(args ...interface{})

	Infof(format string, content ...interface{})
	Info(content ...interface{})

	Errorf(format string, content ...interface{})
	Error(content ...interface{})

	Debugf(format string, content ...interface{})
	Debug(content ...interface{})

	Fatalf(format string, content ...interface{})
	Fatal(content ...interface{})

	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
}

//ILogger 日志接口
type ILogger interface {
	// SetTag(name string, value string)
	ILogging
	GetSessionID() string
	Pause()
	Resume()
}
