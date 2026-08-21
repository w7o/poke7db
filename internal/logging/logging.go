package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/w7o/poke7db/internal/version"
)

func convertDestToLogLocation(dest string) string {
	dest = filepath.Join("logs", dest)
	return dest
}

/*
Writes a message to project/log.txt, non-recoverable
*/
func WriteFile(text string, dest string) {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("\n==\nLOG1: Failed to open or create directory /logs")
	}
	dest = convertDestToLogLocation(dest)

	file, err := os.OpenFile(
		dest,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatalf("\n==\nLOG2: Failure to open file %s", dest)
	}

	defer file.Close()

	if _, err := file.WriteString(text + "\n"); err != nil {
		log.Fatalf("%s\n==\nLOG3: Failure to write message to %s, CHECK IMMEDIATELY\nOFFENDING MESSAGE ABOVE", text, dest)
	}
}

func WriteLog(text string) {
	WriteFile(text, "log.txt")
}

func WriteWarning(text string) {
	WriteLog("WARNING: " + text)
}

func WriteLogAndConsole(text string) {
	log.Println(text)
	WriteLog(text)
}

/*
Returns error object in format:
"code(line num) | message: error"
*/
func RetError(code, message string, err error) error {
	_, _, line, _ := runtime.Caller(1)
	retError := fmt.Errorf("\n%s(%d) | %s: \n\t%w", code, line, message, err)
	if err == nil {
		retError = fmt.Errorf("\n%s(%d) | %s: \n\tnil", code, line, message)
	}
	WriteLog("\n===========\n" + retError.Error())
	return retError
}

/*
Resets the file in ./log/{dest} and creates the file if it doesn't exist;
also creates the log directory if it doesn't exist
*/
func ResetMessageFile(header string, dest string) {
	// create log folder if doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("LOG4: Failed to open or create directory logs")
	}

	dest = convertDestToLogLocation(dest)
	eq := strings.Repeat("=", len(dest)+len(version.G_Version)+3)
	header = fmt.Sprintf("%s [%s]\n%s\n%s\n\n", dest, version.G_Version, eq, header)

	err = os.WriteFile(dest, []byte(header), 0766)
	if err != nil {
		log.Fatalf("LOG5: Failed to initialize %s: %v", dest, err)
	}
}

func ResetLogFile() {
	header := "THIS FILE PROVIDES CONTEXT TO MESSAGES & ERRORS DISPLAYED IN THE CONSOLE"
	ResetMessageFile(header, "log.txt")
}
