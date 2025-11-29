package logger

import (
	"fmt"
	"os"
	"time"
)

type LogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

var logChan = make(chan LogEntry, 1000)

func init() {
	go logWriter()
}

func logWriter() {
	for entry := range logChan {
		fmt.Printf("[%s] %s %s\n", entry.Level, entry.Time.Format("2006-01-02 15:04:05"), entry.Message)
	}
}

func Info(message string) {
	select {
	case logChan <- LogEntry{Level: "INFO", Message: message, Time: time.Now()}:
	default:
		fmt.Fprintln(os.Stderr, "Logger queue is full, dropping log entry")
	}
}

func Error(message string) {
	select {
	case logChan <- LogEntry{Level: "ERROR", Message: message, Time: time.Now()}:
	default:
		fmt.Fprintln(os.Stderr, "Logger queue is full, dropping log entry")
	}
}

func Close() {
	close(logChan)
}
