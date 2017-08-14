package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Log struct {
	Info  *log.Logger
	Debug *log.Logger
}

func (l *Log) InitLogger(logDir, filename string) error {

	folderPath := logDir
	err := os.MkdirAll(folderPath, 0777)
	if nil != err {
		s := fmt.Sprintf("Init logger fail, %s %s", folderPath, err.Error())
		return errors.New(s)
	}

	infoLogName := fmt.Sprintf("%s/%s-info.log", folderPath, filename)
	debugLogName := fmt.Sprintf("%s/%s-debug.log", folderPath, filename)

	infoHandler, err := os.OpenFile(infoLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		s := fmt.Sprintf("Init info log fail, %s %s", infoLogName, err.Error())
		return errors.New(s)
	}

	debugHandler, err := os.OpenFile(debugLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		s := fmt.Sprintf("Init debug log fail, %s %s", debugLogName, err.Error())
		return errors.New(s)
	}

	l.Info = log.New(infoHandler,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Debug = log.New(debugHandler,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
