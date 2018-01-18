package log

import (
	"io"
	"log"
	"os"
)

type BaseLogger struct {
	level int
	out   io.Writer
	log   *log.Logger
}

func NewStdLogger() Logger {
	logger := &BaseLogger{}
	logger.out = os.Stdout
	logger.log = log.New(logger.out, "", log.Ldate|log.Ltime)
	return logger
}

//LOG_DEBUG,LOG_INFO,LOG_WARN,LOG_ERROR
func (logger *BaseLogger) SetLevel(level int) {
	logger.level = level
}

func (logger *BaseLogger) Debug(v ...interface{}) {
	logger.log.Println(v)
}

func (logger *BaseLogger) Debugf(format string, v ...interface{}) {
	logger.log.Printf(format, v...)
}

func (logger *BaseLogger) Info(v ...interface{}) {
	if logger.level > LOG_INFO {
		return
	}
	logger.log.Println(v)
}

func (logger *BaseLogger) Infof(format string, v ...interface{}) {
	if logger.level > LOG_INFO {
		return
	}

	logger.log.Printf(format, v...)
}

func (logger *BaseLogger) Warn(v ...interface{}) {
	if logger.level > LOG_WARN {
		return
	}

	logger.log.Println(v)
}

func (logger *BaseLogger) Warnf(format string, v ...interface{}) {
	if logger.level > LOG_WARN {
		return
	}

	logger.log.Printf(format, v...)
}

func (logger *BaseLogger) Error(v ...interface{}) {
	if logger.level > LOG_ERROR {
		return
	}

	logger.log.Println(v)
}

func (logger *BaseLogger) Errorf(format string, v ...interface{}) {
	if logger.level > LOG_ERROR {
		return
	}

	logger.log.Printf(format, v...)
}

func (logger *BaseLogger) Close() {
	if logger.out != nil {
		//logger.out.Close()
		logger.out = nil
	}
}
