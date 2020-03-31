package log

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var defaultLogger logrusImpl
var defaultLoggerOnce sync.Once

// Data is
type Data struct {
	ClientIP string `` // c.ClientIP() set from controller before entering service
	Session  string `` // ksuid.New().String() or gonanoid.Generate(x, 12) set from first caller
	UserID   string `` // from client
	Type     string `` // MOB/BOF/MSQ/SYS/SCH --> mobile / backoffice / message queuing / system / scheduller
}

// ILogger is
type ILogger interface {
	Debug(sc map[string]interface{}, description string, args ...interface{}) string
	Info(sc map[string]interface{}, description string, args ...interface{}) string
	Warn(sc map[string]interface{}, description string, args ...interface{}) string
	Error(sc map[string]interface{}, description string, args ...interface{}) string
	Fatal(sc map[string]interface{}, description string, args ...interface{}) string
	Panic(sc map[string]interface{}, description string, args ...interface{}) string
	WithFile(path, appsName, serviceName string, maxAge int)
}

// LogrusImpl is
type logrusImpl struct {
	theLogger *logrus.Logger
	useFile   bool
}

// GetLog is
func GetLog() ILogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = logrusImpl{theLogger: logrus.New()}
		defaultLogger.useFile = false
		defaultLogger.theLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "0102 150405.000",
		})
	})
	return &defaultLogger
}

// WithFile is command to state the log will printing to files
// the rolling log file will put in logs/ directory
//
// filename is just a name of log file without any extension
//
// maxAge is age (in days) of the logs file before it gets purged from the file system
func (l *logrusImpl) WithFile(path, appsName, serviceName string, maxAge int) {

	if !l.useFile {

		if maxAge <= 0 {
			panic("maxAge should > 0")
		}

		writer, _ := rotatelogs.New(
			fmt.Sprintf("%s/%s/logs/%s-%s.log.%s", path, appsName, appsName, serviceName, "%Y%m%d"),
			rotatelogs.WithLinkName(fmt.Sprintf("%s/%s/%s.log", path, appsName, serviceName)),
			rotatelogs.WithMaxAge(time.Duration(maxAge*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(1*24)*time.Hour),
		)

		defaultLogger.theLogger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.DebugLevel: writer,
			},
			defaultLogger.theLogger.Formatter,
		))

		l.useFile = true
	}
}

func (l *logrusImpl) getLogEntry(sc map[string]interface{}) *logrus.Entry {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return l.theLogger.WithFields(sc).WithField("func", funcName[x+1:])
}

// Debug is
func (l *logrusImpl) Debug(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Debug(message)
	return message
}

// Info is
func (l *logrusImpl) Info(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Info(message)
	return message
}

// Warn is
func (l *logrusImpl) Warn(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Warn(message)
	return message
}

// Error is
func (l *logrusImpl) Error(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Error(message)
	return message
}

// Fatal is
func (l *logrusImpl) Fatal(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Fatal(message)
	return message
}

// Panic is
func (l *logrusImpl) Panic(sc map[string]interface{}, description string, args ...interface{}) string {
	message := fmt.Sprintf(description, args...)
	l.getLogEntry(sc).Panic(message)
	return message
}
