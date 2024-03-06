// Package app provides functionalities to run the application.
package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Run runs the application.
func Run() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			printError("%T: %s", err, err)
		}
	}()

	if err := rootCommand().Execute(); err != nil {
		printError("%s", err)
	}
}

func printError(format string, args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, color.HiRedString(format, args...))
}
