package logger

import (
	"fmt"
)

// Log simple data to the console
func Log(category string, v ...any) {
	args := append([]any{"[" + category + "] "}, v...)
	fmt.Println(args...)
}

// Logf logs formatted data to the console
func Logf(category, format string, args ...any) {
	format = "[" + category + "] " + format
	fmt.Println(fmt.Sprintf(format, args...))
}
