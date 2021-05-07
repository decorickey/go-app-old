package utils

import (
	"io"
	"log"
	"os"
)

// LoggingSettings ログの基本設定
func LoggingSettings(fileName string) {
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err.Error())
	}
	multiLogfile := io.MultiWriter(os.Stdout, logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogfile)
}
