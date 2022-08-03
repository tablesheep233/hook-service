package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Info.Println("info...")
	go Warning.Println("warning...")
	Error.Println("error...")

	time.Sleep(1 * time.Second)
}
