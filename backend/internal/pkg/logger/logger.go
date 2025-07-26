package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Level 日志级别
type Level int

const (
	// DEBUG 调试级别
	DEBUG Level = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// FATAL 致命级别
	FATAL
)

// Logger 日志记录器
type Logger struct {
	level     Level
	output    io.Writer
	filePath  string
	fileMode  bool
	formatStr string
}

// Config 日志配置
type Config struct {
	Level     string `mapstructure:"level"`
	FilePath  string `mapstructure:"file_path"`
	FileMode  bool   `mapstructure:"file_mode"`
	FormatStr string `mapstructure:"format_str"`
}

// New 创建日志记录器
func New(config Config) (*Logger, error) {
	// 解析日志级别
	var level Level
	switch config.Level {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "warn":
		level = WARN
	case "error":
		level = ERROR
	case "fatal":
		level = FATAL
	default:
		level = INFO
	}

	// 创建日志记录器
	logger := &Logger{
		level:     level,
		output:    os.Stdout,
		filePath:  config.FilePath,
		fileMode:  config.FileMode,
		formatStr: config.FormatStr,
	}

	// 如果启用文件模式，则创建日志文件
	if logger.fileMode {
		if err := logger.createLogFile(); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

// createLogFile 创建日志文件
func (l *Logger) createLogFile() error {
	// 创建日志目录
	dir := filepath.Dir(l.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 打开日志文件
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 设置输出
	l.output = file

	return nil
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", format, args...)
	}
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", format, args...)
	}
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", format, args...)
	}
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", format, args...)
	}
}

// Fatal 记录致命级别日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.level <= FATAL {
		l.log("FATAL", format, args...)
		os.Exit(1)
	}
}

// log 记录日志
func (l *Logger) log(level, format string, args ...interface{}) {
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	// 格式化日志内容
	var content string
	if len(args) > 0 {
		content = fmt.Sprintf(format, args...)
	} else {
		content = format
	}

	// 格式化日志
	var logStr string
	if l.formatStr != "" {
		logStr = fmt.Sprintf(l.formatStr, now, level, content)
	} else {
		logStr = fmt.Sprintf("[%s] [%s] %s\n", now, level, content)
	}

	// 写入日志
	fmt.Fprint(l.output, logStr)
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
	// 如果输出是文件，则关闭文件
	if file, ok := l.output.(*os.File); ok && file != os.Stdout && file != os.Stderr {
		return file.Close()
	}

	return nil
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetOutput 设置日志输出
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// SetFilePath 设置日志文件路径
func (l *Logger) SetFilePath(filePath string) error {
	l.filePath = filePath
	if l.fileMode {
		return l.createLogFile()
	}
	return nil
}

// SetFileMode 设置日志文件模式
func (l *Logger) SetFileMode(fileMode bool) error {
	l.fileMode = fileMode
	if l.fileMode {
		return l.createLogFile()
	}
	return nil
}

// SetFormatStr 设置日志格式
func (l *Logger) SetFormatStr(formatStr string) {
	l.formatStr = formatStr
}