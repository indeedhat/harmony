package logger

import (
	"fmt"
	"log"
	"os"
)

// UseConfigFile when called all log events will be saved to a log file instead of prinnted to the console
func UseConfigFile() error {
	logFile, err := os.OpenFile("harmony.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)
	return nil
}

// Log simple data to the console
func Log(category string, v ...any) {
	args := append([]any{"[" + category + "] "}, v...)
	log.Println(args...)
}

// Logf logs formatted data to the console
func Logf(category, format string, args ...any) {
	format = "[" + category + "] " + format
	log.Println(fmt.Sprintf(format, args...))
}
