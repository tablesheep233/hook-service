package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Fatalln("create log dir fail", err)
	}

	f, err := os.OpenFile("log/hook.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file err", err)
	}

	Info = log.New(io.MultiWriter(f, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(f, os.Stdout), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(f, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
}
