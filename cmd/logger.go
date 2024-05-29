package cmd

import (
    "log"
    "os"
)

var (
    infoLogger  *log.Logger
    errorLogger *log.Logger
)

func init() {
    infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func logInfo(message string) {
    infoLogger.Println(message)
}

func logError(message string) {
    errorLogger.Println(message)
}

