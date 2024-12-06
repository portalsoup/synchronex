package common

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// MultiWriter struct to capture multiple writers
type MultiWriter struct {
	writers []io.Writer
}

// Write writes to all the underlying writers
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func ConfigureLogger(logStdout bool) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	logDir := filepath.Join(usr.HomeDir, ".local", "share", "synchronex", "logs")
	logFile := "synchronex.log"

	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	logPath := filepath.Join(logDir, logFile)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	writers := make([]io.Writer, 0)
	writers = append(writers, file)
	if logStdout {
		writers = append(writers, os.Stdout)
	}

	log.SetOutput(io.MultiWriter(writers...))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	multiWriter := &MultiWriter{
		writers: writers,
	}

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
