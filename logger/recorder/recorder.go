package recorder

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tablesheep233/hook-service/config"
	"github.com/tablesheep233/hook-service/logger"
)

func init() {
	if err := os.MkdirAll(config.Config.ExecLogPath, os.ModePerm); err != nil {
		log.Fatalln("create exec log dir fail", err)
	}
}

type Recorder struct {
	FileName string
	Logger   *log.Logger
}

func NewRecorder(name string) *Recorder {
	fileName := fmt.Sprintf("%s/%s.log", config.Config.ExecLogPath, name)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("open %s file err", name), err)
	}

	return &Recorder{
		FileName: fileName,
		Logger:   log.New(io.MultiWriter(f, os.Stdout), "", log.Ldate|log.Ltime),
	}
}

func (recorder Recorder) Print(msg string) {
	recorder.Logger.Print(msg)
}
