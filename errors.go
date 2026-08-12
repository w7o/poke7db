package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

/*
Returns error object in format:
"code(line num) | message: error"
*/
func retError(code, message string, err error) error {
	_, _, line, _ := runtime.Caller(1)
	if err == nil {
		return fmt.Errorf("%s(%d) | %s: nil", code, line, message)
	}
	return fmt.Errorf("%s(%d) | %s: %w", code, line, message, err)
}

/*
Writes a message to project/error.txt, non-recoverable
*/
func writeError(text string) {
	header := "THIS FILE PROVIDES CONTEXT TO FATAL ERRORS DISPLAYED IN THE CONSOLE\n\n"
	err := os.WriteFile("error.txt", []byte(header+text), 0644)
	if err != nil {
		log.Fatalf("%s\n==\nFailure to write error message, CHECK IMMEDIATELY\nOFFENDING MESSAGE ABOVE", text)
	}
}
