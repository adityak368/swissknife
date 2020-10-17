package logger

import (
	"io"
	"log"
	"os"
)

// Logger Levels
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	file    *os.File = nil
	err     error
)

func initLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// InitConsoleLogger Initializes the console logger
func InitConsoleLogger() {
	initLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

// InitFileLogger Initializes the file logger
func InitFileLogger(filePath string) {
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file ", filePath, ":", err)
	}

	multi := io.MultiWriter(file, os.Stdout)
	multiError := io.MultiWriter(file, os.Stderr)

	initLogger(multi, multi, multi, multiError)
}

// DestroyLogger Cleans Up Logger
func DestroyLogger() {
	if file != nil {
		file.Close()
	}
}
