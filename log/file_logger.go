package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileLogger struct {
	logDir        string
	level         int
	rollingDay    int
	logfile       *os.File
	lastLogFile   *os.File
	stdwriter     io.Writer
	lock          *sync.Mutex
	log           *log.Logger
	logNamePrefix string
}

func NewFileAndStdLogger(logDir string, logNamePrefix string) Logger {
	logger := &FileLogger{}
	logger.logDir = logDir
	logger.logNamePrefix = logNamePrefix
	logger.stdwriter = os.Stdout
	logger.level = LOG_DEBUG
	logger.rollingDay = 0
	logger.lock = new(sync.Mutex)
	return logger
}

func NewFileLogger(logDir string, logNamePrefix string) Logger {
	logger := &FileLogger{}
	logger.logDir = logDir
	logger.logNamePrefix = logNamePrefix
	logger.stdwriter = nil
	logger.level = LOG_DEBUG
	logger.rollingDay = 0
	logger.lock = new(sync.Mutex)
	return logger
}

//LOG_DEBUG,LOG_INFO,LOG_WARN,LOG_ERROR
func (logger *FileLogger) SetLevel(level int) {
	logger.level = level
}

func (logger *FileLogger) Debug(v ...interface{}) {
	logger.rolling()
	logger.log.Println(v)
}

func (logger *FileLogger) Debugf(format string, v ...interface{}) {
	logger.rolling()
	logger.log.Printf(format, v...)
}

func (logger *FileLogger) Info(v ...interface{}) {
	if logger.level > LOG_INFO {
		return
	}
	logger.rolling()
	logger.log.Println(v)
}

func (logger *FileLogger) Infof(format string, v ...interface{}) {
	if logger.level > LOG_INFO {
		return
	}
	logger.rolling()
	logger.log.Printf(format, v...)
}

func (logger *FileLogger) Warn(v ...interface{}) {
	if logger.level > LOG_WARN {
		return
	}
	logger.rolling()
	logger.log.Println(v)
}

func (logger *FileLogger) Warnf(format string, v ...interface{}) {
	if logger.level > LOG_WARN {
		return
	}
	logger.rolling()
	logger.log.Printf(format, v...)
}

func (logger *FileLogger) Error(v ...interface{}) {
	if logger.level > LOG_ERROR {
		return
	}
	logger.rolling()
	logger.log.Println(v)
}

func (logger *FileLogger) Errorf(format string, v ...interface{}) {
	if logger.level > LOG_ERROR {
		return
	}
	logger.rolling()
	logger.log.Printf(format, v...)
}

func (logger *FileLogger) Close() {
	if logger.logfile != nil {
		logger.logfile.Close()
		logger.logfile = nil
	}
}

func (logger *FileLogger) rolling() {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger._rolling()
}

func (logger *FileLogger) _rolling() {
	now := time.Now()
	if logger.rollingDay < 0 {
		if logger.lastLogFile != nil {
			logger.lastLogFile.Close()
			logger.lastLogFile = nil
		}
		return
	} else {
		if logger.logfile != nil {
			logger.logfile.Sync()
		}
		//延迟关闭上一个log文件, 优化锁粒度: 无需锁文件写, 只需锁rolling本身
		logger.lastLogFile = logger.logfile
		logger.rollingDay = now.YearDay()

		logName := logger.logNamePrefix + now.Format("20060102") + ".log"
		logPath := filepath.Join(logger.logDir, logName)

		var err error
		logger.logfile, err = os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("fail to create log file %v err=%v\n", logPath, err)
		}

		if logger.stdwriter != nil {
			logger.log = log.New(io.MultiWriter(logger.logfile, logger.stdwriter), "", log.Ldate|log.Ltime)
		} else {
			logger.log = log.New(logger.logfile, "", log.Ldate|log.Ltime)
		}
	}
}
