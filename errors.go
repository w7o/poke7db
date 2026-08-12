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
	retError := fmt.Errorf("\n%s(%d) | %s: \n\t%w", code, line, message, err)
	if err == nil {
		retError = fmt.Errorf("\n%s(%d) | %s: \n\tnil", code, line, message)
	}
	writeLog("\n===========\n" + retError.Error())
	return retError
}

/*
Writes a message to project/log.txt, non-recoverable
*/
func writeLog(text string) {
	file, err := os.OpenFile(
		"log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal("\n==\nFailure to open file")
	}

	defer file.Close()

	if _, err := file.WriteString(text + "\n"); err != nil {
		log.Fatalf("%s\n==\nFailure to write log message, CHECK IMMEDIATELY\nOFFENDING MESSAGE ABOVE", text)
	}
}

func clearErrorFile() {
	header := "THIS FILE PROVIDES CONTEXT TO MESSAGES & ERRORS DISPLAYED IN THE CONSOLE\n\n"

	err := os.WriteFile("log.txt", []byte(header), 0666)
	if err != nil {
		log.Fatalf("Failed to initialize log.txt: %v", err)
	}
}
