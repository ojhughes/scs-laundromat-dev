package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	log.SetOutput(os.Stdout)
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}
