package tteok

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

var out io.Writer = os.Stdout

type log struct {
	Process   string                 `json:"process"`
	Level     string                 `json:"level"`
	Timestamp string                 `json:"timestamp"` //RFC3339
	At        string                 `json:"at"`
	Message   string                 `json:"message,omitempty"`
	Messages  []string               `json:"messages,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
	Func      string                 `json:"func,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Debug outputs a standarized Tteok log message with level "DEBUG"
func Debug(params ...interface{}) { tteok("DEBUG", params) }

// Info outputs a standarized Tteok log message with level "INFO"
func Info(params ...interface{}) { tteok("INFO", params) }

// Warn outputs a standarized Tteok log message with level "WARN"
func Warn(params ...interface{}) { tteok("WARN", params) }

// Fatal outputs a standarized Tteok log message with level "FATAL"
func Fatal(params ...interface{}) { tteok("FATAL", params) }

func tteok(level string, params []interface{}) {
	logMsg := buildLog(level, params)

	logMsgJSON, err := json.MarshalIndent(logMsg, "", "  ")
	if err != nil {
		return
	}

	if level == "DEBUG" && os.Getenv("DEBUG") != "true" {
		return
	}

	fmt.Fprintf(out, "%+v\n", string(logMsgJSON))
}

func buildLog(level string, params []interface{}) log {
	logMsg := log{
		Process:   os.Args[0],
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	logMsg.addCallingLocation()

	for _, param := range params {
		logMsg.enrich(param)
	}

	return logMsg
}

func (l *log) addCallingLocation() {
	defer func() {
		if e := recover(); e != nil {
			stack := string(debug.Stack())
			stackSplit := strings.Split(stack, "\n")

			callingLine := stackSplit[16] // always fixed (FIFO)
			callingLineSplit := strings.Split(callingLine, "/")

			callingLocation := callingLineSplit[len(callingLineSplit)-1]
			callingLocation = strings.Split(callingLocation, " ")[0] // format: "\tfile\path\here\file.go:123 0xfff"

			l.At = callingLocation
		}
	}()

	panic(fmt.Errorf(""))
}

// enrich determines the type of the incoming data and dispatches accordingly
// used recursively for arrays
func (l *log) enrich(v interface{}) {
	switch v.(type) {
	case error:
		l.addError(v)
	case string:
		l.addMessage(v)
	default:
		switch reflect.ValueOf(v).Kind() {
		case reflect.Ptr:
			l.dereferencePointer(v)
		case reflect.Map:
			l.addMapToData(v)
		case reflect.Slice:
			l.spreadSlice(v)
		case reflect.Struct:
			l.addStructToData(v)
		default:
			// we ignore numbers, funcs,
		}
	}
}
