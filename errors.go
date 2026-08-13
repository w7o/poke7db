package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
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

func writeFile(text string, dest string) {
	file, err := os.OpenFile(
		dest,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatalf("\n==\nFailure to open file %s", dest)
	}

	defer file.Close()

	if _, err := file.WriteString(text + "\n"); err != nil {
		log.Fatalf("%s\n==\nFailure to write message to %s, CHECK IMMEDIATELY\nOFFENDING MESSAGE ABOVE", text, dest)
	}
}

func writeLog(text string) {
	writeFile(text, "log.txt")
}

func writeWarning(text string) {
	writeLog("WARNING: " + text)
}

func writeLogAndConsole(text string) {
	log.Println(text)
	writeLog(text)
}

func resetMessageFile(header string, dest string) {
	eq := strings.Repeat("=", len(dest)+len(VERSION)+3)
	header = fmt.Sprintf("%s [%s]\n%s\n%s\n\n", dest, VERSION, eq, header)
	err := os.WriteFile(dest, []byte(header), 0666)
	if err != nil {
		log.Fatalf("Failed to initialize %s: %v", dest, err)
	}
}

func resetLogFile() {
	header := "THIS FILE PROVIDES CONTEXT TO MESSAGES & ERRORS DISPLAYED IN THE CONSOLE"
	resetMessageFile(header, "log.txt")
}
