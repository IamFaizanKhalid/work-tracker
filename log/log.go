package log

import (
	"github.com/IamFaizanKhalid/work-tracker/constants"
	"io"
	"log"
	"os"
)

var (
	Info       *log.Logger
	Warn       *log.Logger
	Error      *log.Logger
	OutputFile *os.File
)

func InitLogger() {
	var err error
	OutputFile, err = os.OpenFile(constants.WorkDir+"/"+constants.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defaultWriter := io.Writer(OutputFile)

	Info = log.New(defaultWriter, "[INFO]\t", log.LstdFlags|log.Lshortfile)
	Warn = log.New(defaultWriter, "[WARN]\t", log.LstdFlags|log.Lshortfile)
	Error = log.New(defaultWriter, "[ERROR]\t", log.LstdFlags|log.Lshortfile)
}
