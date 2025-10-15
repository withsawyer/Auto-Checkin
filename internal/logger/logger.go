package logger

import (
	"auto-checkin/internal/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// 日志级别
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var (
	instance *Logger
	once     sync.Once
)

type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
	mu          sync.Mutex
	level       int
	maxSize     int64
	maxBackups  int
}

// GetLogger 获取日志实例(单例模式)
func Log() *Logger {
	once.Do(func() {
		instance = &Logger{
			level:      DEBUG,
			maxSize:    10 * 1024 * 1024, // 10MB
			maxBackups: 5,
		}
		if config.Cfg.Debug {
			instance.debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		}
	})
	return instance
}

// Init 初始化日志
func (l *Logger) Init(filename string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 创建日志目录
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	// 检查是否需要轮转
	if l.shouldRotate(filename) {
		if err := l.rotate(filename); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = file

	// 动态初始化日志输出目标
	if config.Cfg.Debug {
		// Debug 模式下，所有日志等级输出到控制台和文件
		l.debugLogger = log.New(io.MultiWriter(file, os.Stdout), "DEBUG: ", log.Ldate|log.Ltime)
		l.infoLogger = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime)
		l.warnLogger = log.New(io.MultiWriter(file, os.Stdout), "WARN: ", log.Ldate|log.Ltime)
		l.errorLogger = log.New(io.MultiWriter(file, os.Stdout), "ERROR: ", log.Ldate|log.Ltime)
	} else {
		// 非 Debug 模式下，仅输出到文件
		l.debugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime)
		l.infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
		l.warnLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime)
		l.errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	}

	return nil
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// Debug 记录调试信息
func (l *Logger) Debug(v ...interface{}) {
	if config.Cfg.Debug && l.level <= DEBUG && l.debugLogger != nil {
		l.debugLogger.Println(v...)
	}
}

// Debugf 格式化记录调试信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	if config.Cfg.Debug && l.level <= DEBUG && l.debugLogger != nil {
		l.debugLogger.Printf(format, v...)
	}
}

// Info 记录普通信息
func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO && l.infoLogger != nil {
		l.infoLogger.Println(v...)
	}
}

// Infof 格式化记录普通信息
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO && l.infoLogger != nil {
		l.infoLogger.Printf(format, v...)
	}
}

// Warn 记录警告信息
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN && l.warnLogger != nil {
		l.warnLogger.Println(v...)
	}
}

// Warn 记录警告信息
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN && l.warnLogger != nil {
		l.warnLogger.Printf(format, v...)
	}
}

// Error 记录错误信息
func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR && l.errorLogger != nil {
		l.errorLogger.Println(v...)
	}
}

// Errorf 格式化记录错误信息
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR && l.errorLogger != nil {
		l.errorLogger.Printf(format, v...)
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// shouldRotate 检查是否需要轮转
func (l *Logger) shouldRotate(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.Size() >= l.maxSize
}

// rotate 执行日志轮转
func (l *Logger) rotate(filename string) error {
	// 关闭当前日志文件
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return err
		}
	}

	// 归档旧日志文件
	for i := l.maxBackups - 1; i >= 1; i-- {
		oldName := filename + "." + time.Now().Format("20060102") + "." + strconv.Itoa(i)
		newName := filename + "." + time.Now().Format("20060102") + "." + strconv.Itoa(i+1)
		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	// 重命名当前日志文件
	if _, err := os.Stat(filename); err == nil {
		newName := filename + "." + time.Now().Format("20060102") + ".1"
		if err := os.Rename(filename, newName); err != nil {
			return err
		}
	}

	return nil
}
