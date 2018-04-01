// Package logx is simple extended version of standard 'log' package based on
// logLevel. Most of the concepts inspired from https://godoc.org/github.com/golang/glog
// and https://github.com/goinggo/tracelog. These packages are huge and complex.
// Hence we writing our own log package with as simple as possible.
// Another ref: https://github.com/UlricQin/goutils/blob/master/logtool/logtool.go
package logx

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// logLevel is a severity level at which logger works.
type logLevel int

const (
	// InfoLevel ...
	InfoLevel logLevel = iota
	// WarningLevel ...
	WarningLevel
	// ErrorLevel ...
	ErrorLevel
	// FatalLevel ...
	FatalLevel
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger

	// Levels ...
	Levels = map[string]logLevel{
		"INFO":    InfoLevel,
		"WARNING": WarningLevel,
		"ERROR":   ErrorLevel,
		"FATAL":   FatalLevel,
	}
)

func init() {
	Init(InfoLevel, nil)
}

// Init ...
func Init(lev logLevel, multiHandler io.Writer) {
	infoHandler := ioutil.Discard
	warningHandler := ioutil.Discard
	errorHandler := ioutil.Discard
	fatalHandler := ioutil.Discard

	switch lev {
	case InfoLevel:
		infoHandler = os.Stdout
		warningHandler = os.Stdout
		errorHandler = os.Stderr
		fatalHandler = os.Stderr
	case WarningLevel:
		warningHandler = os.Stdout
		errorHandler = os.Stderr
		fatalHandler = os.Stderr
	case ErrorLevel:
		errorHandler = os.Stderr
		fatalHandler = os.Stderr
	case FatalLevel:
		fatalHandler = os.Stderr
	default:
		log.Fatal("logx: Invalid log level should be (0-3)")
	}

	if multiHandler != nil {
		if infoHandler == os.Stdout {
			infoHandler = io.MultiWriter(infoHandler, multiHandler)
		}
		if warningHandler == os.Stdout {
			warningHandler = io.MultiWriter(warningHandler, multiHandler)
		}
		if errorHandler == os.Stderr {
			errorHandler = io.MultiWriter(errorHandler, multiHandler)
		}
		if fatalHandler == os.Stderr {
			fatalHandler = io.MultiWriter(fatalHandler, multiHandler)
		}
	}

	infoLogger = log.New(infoHandler, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(warningHandler, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errorHandler, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(fatalHandler, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// -----------------------------------------------------------------------------
// Wrapper for changing logger's output writer anytime.

// InfoSetOutput ...
func InfoSetOutput(w io.Writer) {
	infoLogger.SetOutput(w)
}

// WarningSetOutput ...
func WarningSetOutput(w io.Writer) {
	warningLogger.SetOutput(w)
}

// ErrorSetOutput ...
func ErrorSetOutput(w io.Writer) {
	errorLogger.SetOutput(w)
}

// FatalSetOutput ...
func FatalSetOutput(w io.Writer) {
	fatalLogger.SetOutput(w)
}

// -----------------------------------------------------------------------------

// Info writes into infologger as same as basic fmt.Print.
func Info(v ...interface{}) {
	infoLogger.Output(2, fmt.Sprint(v...))
}

// Infoln writes into infologger as same as basic fmt.Println.
func Infoln(v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintln(v...))
}

// Infof writes into infologger as same as basic Outputf.
func Infof(format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(format, v...))
}

// Warning writes warning messages into warninglogger as same as basic
// log.Output.
func Warning(v ...interface{}) {
	warningLogger.Output(2, fmt.Sprint(v...))
}

// Warningln writes warning messages into warninglogger as same as basic
// log.Outputln.
func Warningln(v ...interface{}) {
	warningLogger.Output(2, fmt.Sprintln(v...))
}

// Warningf writes warning messages into warninglogger as same as basic
// log.Outputf.
func Warningf(format string, v ...interface{}) {
	warningLogger.Output(2, fmt.Sprintf(format, v...))
}

// Error writes error messages into errorlogger as same as basic log.Error.
func Error(v ...interface{}) {
	errorLogger.Output(2, fmt.Sprint(v...))
}

// Errorln writes error messages into errorlogger as same as basic log.Errorln.
func Errorln(v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintln(v...))
}

// Errorf writes error messages into errorlogger as same as basic log.Errorf.
func Errorf(format string, v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}

// Fatal writes fatal error messages into errorlogger and exit as same as basic
// log.Fatal.
func Fatal(v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalln writes fatal error messages into errorlogger and exit as same as
// basic log.Fataln.
func Fatalln(v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf writes fatal error messages into errorlogger and exit as same as basic
// log.Fataf.
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// -----------------------------------------------------------------------------
