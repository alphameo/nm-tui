// Package logger provides simple tools for debugging, system messaging and other
package logger

import (
	"log"
	"os"
)

const LoggerFlags = log.Ldate | log.Ltime | log.Lshortfile

type LoggerLevel int

const (
	InformationLevel LoggerLevel = iota
	WarningsLevel
	ErrorsLevel
)

var (
	infoLog  = log.New(os.Stdout, "INFO: ", LoggerFlags)
	warnLog  = log.New(os.Stdout, "WARNING: ", LoggerFlags)
	errLog   = log.New(os.Stderr, "ERROR: ", LoggerFlags)
	debugLog = log.New(os.Stdout, "DEBUG: ", LoggerFlags)
	Level    = ErrorsLevel
)

func FilePath(path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	infoLog.SetOutput(f)
	warnLog.SetOutput(f)
	errLog.SetOutput(f)
	debugLog.SetOutput(f)
}

func Inform(v ...any) {
	if Level > InformationLevel {
		return
	}
	infoLog.Print(v...)
}

func Informln(v ...any) {
	if Level > InformationLevel {
		return
	}
	infoLog.Println(v...)
}

func Informf(format string, v ...any) {
	if Level > InformationLevel {
		return
	}
	infoLog.Printf(format, v...)
}

func Warn(v ...any) {
	if Level > WarningsLevel {
		return
	}
	warnLog.Print(v...)
}

func Warnln(v ...any) {
	if Level > WarningsLevel {
		return
	}
	warnLog.Println(v...)
}

func Warnf(format string, v ...any) {
	if Level > WarningsLevel {
		return
	}
	warnLog.Printf(format, v...)
}

func Err(v ...any) {
	errLog.Print(v...)
}

func Errln(v ...any) {
	errLog.Println(v...)
}

func Errf(format string, v ...any) {
	errLog.Printf(format, v...)
}

func Debug(v ...any) {
	debugLog.Print(v...)
}

func Debugln(v ...any) {
	debugLog.Println(v...)
}

func Debugf(format string, v ...any) {
	debugLog.Printf(format, v...)
}
