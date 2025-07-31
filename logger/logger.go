// Package logger provides a logging utility for Go applications.
package logger

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	manifest "github.com/rafa-mori/goforge/info"
	l "github.com/rafa-mori/logz"
)

type GLog[T any] interface {
	GetLogger() l.Logger
	GetLogLevel() LogLevel
	GetShowTrace() bool
	GetDebug() bool
	SetLogLevel(string)
	SetDebug(bool)
	SetShowTrace(bool)
	ObjLog(*T, string, ...string)
	Log(string, ...any)
}
type gLog[T any] struct {
	l.Logger
	gLogLevel  LogLevel // Global log level
	gShowTrace bool     // Flag to show trace in logs
	gDebug     bool     // Flag to show debug messages
}
type LogType string
type LogLevel int

var (
	info      manifest.Manifest
	debug     bool
	showTrace bool
	logLevel  string
	g         *gLog[l.Logger] // Global logger instance
	Logger    GLog[l.Logger]
	err       error
)

const (
	// LogTypeDebug is the log type for debug messages.
	LogTypeDebug LogType = "debug"
	// LogTypeNotice is the log type for notice messages.
	LogTypeNotice LogType = "notice"
	// LogTypeInfo is the log type for informational messages.
	LogTypeInfo LogType = "info"
	// LogTypeWarn is the log type for warning messages.
	LogTypeWarn LogType = "warn"
	// LogTypeError is the log type for error messages.
	LogTypeError LogType = "error"
	// LogTypeFatal is the log type for fatal error messages.
	LogTypeFatal LogType = "fatal"
	// LogTypePanic is the log type for panic messages.
	LogTypePanic LogType = "panic"
	// LogTypeSuccess is the log type for success messages.
	LogTypeSuccess LogType = "success"
)

const (
	// LogLevelDebug is the log level for debug messages.
	LogLevelDebug LogLevel = iota
	// LogLevelNotice is the log level for notice messages.
	LogLevelNotice
	// LogLevelInfo is the log level for informational messages.
	LogLevelInfo
	// LogLevelSuccess is the log level for success messages.
	LogLevelSuccess
	// LogLevelWarn is the log level for warning messages.
	LogLevelWarn
	// LogLevelError is the log level for error messages.
	LogLevelError
	// LogLevelFatal is the log level for fatal error messages.
	LogLevelFatal
	// LogLevelPanic is the log level for panic messages.
	LogLevelPanic
)

func getEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	} else {
		valInterface := reflect.ValueOf(value)
		if valInterface.Type().ConvertibleTo(reflect.TypeFor[T]()) {
			return valInterface.Convert(reflect.TypeFor[T]()).Interface().(T)
		}
	}
	return defaultValue
}

func init() {
	if info == nil {
		info, err = manifest.GetManifest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get info manifest: %v\n", err)
			os.Exit(1)
		}
		l.GetLogger(info.GetBin())
	}
	if Logger == nil {
		Logger = GetLogger[l.Logger](nil)
		if logger, ok := Logger.(*gLog[l.Logger]); ok {
			g = logger
			logLevel = getEnvOrDefault("GOBE_LOG_LEVEL", "error")
			debug = getEnvOrDefault("GOBE_DEBUG", false)
			showTrace = getEnvOrDefault("GOBE_SHOW_TRACE", false)
			g.gLogLevel = LogLevelError
			g.gShowTrace = showTrace
			g.gDebug = debug
		}
	}
}

func SetDebug(d bool) {
	if g == nil || Logger == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.gDebug = d
	if d {
		g.SetLevel("debug")
	} else {
		switch g.gLogLevel {
		case LogLevelDebug:
			g.SetLevel("debug")
		case LogLevelInfo:
			g.SetLevel("info")
		case LogLevelWarn:
			g.SetLevel("warn")
		case LogLevelError:
			g.SetLevel("error")
		case LogLevelFatal:
			g.SetLevel("fatal")
		case LogLevelPanic:
			g.SetLevel("panic")
		case LogLevelNotice:
			g.SetLevel("notice")
		case LogLevelSuccess:
			g.SetLevel("success")
		default:
			g.SetLevel("info")
		}
	}
}
func setLogLevel(logLevel string) {
	if g == nil || Logger == nil {
		_ = GetLogger[l.Logger](nil)
	}
	switch strings.ToLower(logLevel) {
	case "debug":
		g.gLogLevel = LogLevelDebug
		g.SetLevel("debug")
	case "info":
		g.gLogLevel = LogLevelInfo
		g.SetLevel("info")
	case "warn":
		g.gLogLevel = LogLevelWarn
		g.SetLevel("warn")
	case "error":
		g.gLogLevel = LogLevelError
		g.SetLevel("error")
	case "fatal":
		g.gLogLevel = LogLevelFatal
		g.SetLevel("fatal")
	case "panic":
		g.gLogLevel = LogLevelPanic
		g.SetLevel("panic")
	case "notice":
		g.gLogLevel = LogLevelNotice
		g.SetLevel("notice")
	case "success":
		g.gLogLevel = LogLevelSuccess
		g.SetLevel("success")
	default:
		logLevel = "error"
		g.gLogLevel = LogLevelError
		g.SetLevel(logLevel)
	}
}
func getShowTrace() bool {
	if debug {
		return true
	} else {
		if !showTrace {
			return false
		} else {
			return true
		}
	}
}
func willPrintLog(logType string) bool {
	if debug {
		return true
	} else {
		lTypeInt := LogLevelError
		switch strings.ToLower(logType) {
		case "debug":
			lTypeInt = 0
		case "fatal":
			lTypeInt = 0
		case "panic":
			lTypeInt = 0
		case "info":
			lTypeInt = LogLevelInfo
		case "warn":
			lTypeInt = LogLevelWarn
		case "error":
			lTypeInt = LogLevelError
		case "notice":
			lTypeInt = LogLevelNotice
		case "success":
			lTypeInt = LogLevelSuccess
		default:
			lTypeInt = LogLevelError
		}
		return lTypeInt >= g.gLogLevel
	}
}
func GetLogger[T any](obj *T) GLog[l.Logger] {
	if g == nil || Logger == nil {
		g = &gLog[l.Logger]{
			Logger:     l.GetLogger(info.GetBin()),
			gLogLevel:  LogLevelInfo,
			gShowTrace: showTrace,
			gDebug:     debug,
		}
		Logger = g
	}
	if obj == nil {
		return Logger
	}
	var lgr l.Logger
	if objValueLogger := reflect.ValueOf(obj).Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
		if objValueLogger = reflect.ValueOf(obj).Elem().FieldByName("Logger"); !objValueLogger.IsValid() {
			g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  "error",
				"object":   obj,
				"msg":      "object does not have a logger field",
				"showData": getShowTrace(),
			})
			return g
		} else {
			lgrC := objValueLogger.Convert(reflect.TypeFor[l.Logger]())
			if lgrC.IsNil() {
				lgrC = reflect.ValueOf(g.Logger)
			}
			if lgr = lgrC.Interface().(l.Logger); lgr == nil {
				lgr = g.Logger
			}
		}
	} else {
		lgr = g
	}
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  "error",
			"object":   obj,
			"msg":      "object does not have a logger field",
			"showData": getShowTrace(),
		})
		return Logger
	}
	return &gLog[l.Logger]{
		Logger:     lgr,
		gLogLevel:  g.gLogLevel,
		gShowTrace: g.gShowTrace,
		gDebug:     g.gDebug,
	}
}
func getCtxMessageMap(logType, funcName, file string, line int) map[string]any {
	ctxMessageMap := map[string]any{
		"context":   funcName,
		"file":      file,
		"line":      line,
		"logType":   logType,
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   info.GetVersion(),
	}
	if !debug && !showTrace {
		ctxMessageMap["showData"] = false
	} else {
		ctxMessageMap["showData"] = getShowTrace()
	}
	if info != nil {
		ctxMessageMap["appName"] = info.GetName()
		ctxMessageMap["bin"] = info.GetBin()
		ctxMessageMap["version"] = info.GetVersion()
	}
	return ctxMessageMap
}
func LogObjLogger[T any](obj *T, logType string, messages ...string) {
	lgr := GetLogger(obj)
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": getShowTrace(),
		})
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		lgr.GetLogger().ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fullMessage := strings.Join(messages, " ")
	logType = strings.ToLower(logType)

	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			logging(lgr.GetLogger(), lType, fullMessage, ctxMessageMap)
		} else {
			lgr.GetLogger().ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		lgr.GetLogger().InfoCtx(fullMessage, ctxMessageMap)
	}
}
func Log(logType string, messages ...any) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		g.ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fullMessage := ""
	if len(messages) > 0 {
		fullMessage = fmt.Sprintf("%v", messages[0:])
	}
	logType = strings.ToLower(logType)
	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g.Logger, lType, fullMessage, ctxMessageMap)
		} else {
			g.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		logging(g.Logger, LogTypeInfo, fullMessage, ctxMessageMap)
	}
}
func logging(lgr l.Logger, lType LogType, fullMessage string, ctxMessageMap map[string]any) {
	lt := strings.ToLower(string(lType))
	if _, exist := ctxMessageMap["showData"]; !exist {
		ctxMessageMap["showData"] = getShowTrace()
	}
	if willPrintLog(lt) {
		switch lType {
		case LogTypeInfo:
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		case LogTypeDebug:
			lgr.DebugCtx(fullMessage, ctxMessageMap)
		case LogTypeError:
			lgr.ErrorCtx(fullMessage, ctxMessageMap)
		case LogTypeWarn:
			lgr.WarnCtx(fullMessage, ctxMessageMap)
		case LogTypeNotice:
			lgr.NoticeCtx(fullMessage, ctxMessageMap)
		case LogTypeSuccess:
			lgr.SuccessCtx(fullMessage, ctxMessageMap)
		case LogTypeFatal:
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		case LogTypePanic:
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		default:
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		}
	} else {
		ctxMessageMap["msg"] = fullMessage
		ctxMessageMap["showData"] = false
		lgr.DebugCtx("Log: message not printed due to log level", ctxMessageMap)
	}
}

func (g *gLog[T]) GetLogger() l.Logger                 { return g.Logger }
func (g *gLog[T]) GetLogLevel() LogLevel               { return g.gLogLevel }
func (g *gLog[T]) GetShowTrace() bool                  { return g.gShowTrace }
func (g *gLog[T]) GetDebug() bool                      { return g.gDebug }
func (g *gLog[T]) SetLogLevel(logLevel string)         { setLogLevel(logLevel) }
func (g *gLog[T]) SetShowTrace(showTrace bool)         { g.gShowTrace = showTrace }
func (g *gLog[T]) SetDebug(d bool)                     { SetDebug(d); g.gDebug = d }
func (g *gLog[T]) Log(logType string, messages ...any) { Log(logType, messages...) }
func (g *gLog[T]) ObjLog(obj *T, logType string, messages ...string) {
	LogObjLogger(obj, logType, messages...)
}

func NewLogger[T any](prefix string) GLog[T] {
	return &gLog[T]{
		Logger:     l.NewLogger(prefix),
		gLogLevel:  LogLevelError,
		gShowTrace: false,
		gDebug:     false,
	}
}
