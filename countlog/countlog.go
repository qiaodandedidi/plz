package countlog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// push event out

const LevelTrace = 10
const LevelDebug = 20
const LevelInfo = 30
const LevelWarn = 40
const LevelError = 50
const LevelFatal = 60

// MinLevel exists to minimize the overhead of Trace/Debug logging
var MinLevel = LevelDebug

func getLevelName(level int) string {
	switch level {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func ShouldLog(level int) bool {
	return level >= MinLevel
}

func TraceCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		logCall(LevelError, callee, err, properties)
		return
	}
	if LevelTrace < MinLevel {
		return
	}
	logCall(LevelTrace, callee, err, properties)
}

//go:noinline
func logCall(level int, callee string, err error, properties []interface{}) {
	callee = callee[len("callee!"):]
	log(level, "event!call "+callee, append(properties, "callee", callee, "err", err))
}

func Debug(event string, properties ...interface{}) {
	if LevelDebug < MinLevel {
		return
	}
	log(LevelDebug, event, properties)
}

func DebugCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		logCall(LevelError, callee, err, properties)
		return
	}
	if LevelDebug < MinLevel {
		return
	}
	logCall(LevelDebug, callee, err, properties)
}

func Info(event string, properties ...interface{}) {
	if LevelInfo < MinLevel {
		return
	}
	log(LevelInfo, event, properties)
}

func InfoCall(callee string, err error, properties ...interface{}) {
	if err != nil {
		logCall(LevelError, callee, err, properties)
		return
	}
	if LevelInfo < MinLevel {
		return
	}
	logCall(LevelInfo, callee, err, properties)
}

func Warn(event string, properties ...interface{}) {
	log(LevelWarn, event, properties)
}

func Error(event string, properties ...interface{}) {
	log(LevelError, event, properties)
}

func Fatal(event string, properties ...interface{}) {
	log(LevelFatal, event, properties)
}

func Log(level int, event string, properties ...interface{}) {
	log(level, event, properties)
}

func log(level int, event string, properties []interface{}) {
	var expandedProperties []interface{}
	if len(LogWriters) == 0 {
		if expandedProperties == nil {
			level, event, expandedProperties = expand(level, event, properties)
			if level < MinLevel {
				return
			}
		}
		defaultLogWriter.WriteLog(level, event, expandedProperties)
		return
	}
	for _, logWriter := range LogWriters {
		if !logWriter.ShouldLog(level, event, properties) {
			continue
		}
		if expandedProperties == nil {
			level, event, expandedProperties = expand(level, event, properties)
			if level < MinLevel {
				return
			}
		}
		logWriter.WriteLog(level, event, expandedProperties)
	}
}
func expand(level int, event string, properties []interface{}) (int, string, []interface{}) {
	if strings.HasPrefix(event, "event!") {
		event = event[len("event!"):]
	} else {
		_, file, line, _ := runtime.Caller(3)
		// this format allows intellij to jump to that line
		lineNumber := fmt.Sprintf("%s:%d", file, line)
		os.Stderr.Write([]byte("countlog event must starts with event! prefix:" + lineNumber + "\n"))
		os.Stderr.Sync()
	}
	expandedProperties := make([]interface{}, 0, len(properties)+4)
	expandedProperties = append(expandedProperties, "timestamp")
	expandedProperties = append(expandedProperties, time.Now().UnixNano())
	skipFramesCount := 3
	for _, prop := range properties {
		switch typedProp := prop.(type) {
		case func() interface{}:
			expandedProperties = append(expandedProperties, typedProp())
		case []byte:
			// []byte is likely being reused, need to make a copy here
			expandedProperties = append(expandedProperties, encodeAnyByteArray(typedProp))
		case Context:
			skipFramesCount = 5
			expandedProperties = append(expandedProperties, prop)
		default:
			expandedProperties = append(expandedProperties, prop)
		}
	}
	_, file, line, ok := runtime.Caller(skipFramesCount)
	if ok {
		expandedProperties = append(expandedProperties, "lineNumber")
		// this format allows intellij to jump to that line
		lineNumber := fmt.Sprintf("%s:%d", file, line)
		expandedProperties = append(expandedProperties, lineNumber)
	}
	return level, event, expandedProperties
}

// pull state callbacks

// like JMX MBean
type StateExporter interface {
	ExportState() map[string]interface{}
}

var stateExporters = map[string]StateExporter{}
var stateExportersMutex = &sync.Mutex{}

func RegisterStateExporter(name string, se StateExporter) {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	stateExporters[name] = se
}

func RegisterStateExporterByFunc(name string, f func() map[string]interface{}) {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	stateExporters[name] = &funcStateExporter{f}
}

func StateExporters() map[string]StateExporter {
	stateExportersMutex.Lock()
	defer stateExportersMutex.Unlock()
	return stateExporters
}

type funcStateExporter struct {
	f func() map[string]interface{}
}

func (se *funcStateExporter) ExportState() map[string]interface{} {
	return se.f()
}
