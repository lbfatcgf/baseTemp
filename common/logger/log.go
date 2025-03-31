package logger

import (
	"baseTemp/common/config"
	"io"
	"log/slog"
	"os"
	fp "path/filepath"
	"time"
)

var (
	logger     *slog.Logger
	logMsgChan chan *logMsg
	logFile    *os.File

)

type logMsg struct {
	Level   slog.Level
	Message string
	arg     []any
}

func LogInfo(msg string, arg ...any) {

	logMsgChan <- &logMsg{Level: slog.LevelInfo, Message: msg, arg: arg}
}

func LogError(msg string, arg ...any) {
	logMsgChan <- &logMsg{Level: slog.LevelError, Message: msg, arg: arg}
}

func LogWarn(msg string, arg ...any) {
	logMsgChan <- &logMsg{Level: slog.LevelWarn, Message: msg, arg: arg}
}
func LogDebug(msg string, arg ...any) {
	logMsgChan <- &logMsg{Level: slog.LevelDebug, Message: msg, arg: arg}
}
func InitLog(mode string) {
	logMsgChan = make(chan *logMsg, 1000)
	logdir := config.Conf().LogDir
	// fmt.Println(fp.Join(logdir,Conf().Name+"-"+time.Now().Format("2006-01-02")+".log"))
	logFile, err := os.OpenFile(fp.Join(logdir, config.Conf().Name+"-"+time.Now().Format("2006-01-02")+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)

	}
	if mode == "release" {

		logHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(logHandler)
	} else {
		writer := io.MultiWriter(os.Stdout, logFile)
		logHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(logHandler)
	}
	go func() {
		for {
			msg := <-logMsgChan
			if msg == nil {
				continue
			}
			checkLogFileName()
			switch msg.Level {
			case slog.LevelDebug:
				logger.Debug(msg.Message, msg.arg...)
			case slog.LevelInfo:
				logger.Info(msg.Message, msg.arg...)
			case slog.LevelError:
				logger.Error(msg.Message, msg.arg...)
			case slog.LevelWarn:
				logger.Warn(msg.Message, msg.arg...)
			}
		}
	}()

}

func checkLogFileName() {
	logdir := config.Conf().LogDir
	newFilePtah:=fp.Join(logdir, config.Conf().Name+"-"+time.Now().Format("2006-01-02")+".log")
	if newFilePtah!=logFile.Name(){
		logFile.Close()
		nf, err := os.OpenFile(fp.Join(config.Conf().LogDir, newFilePtah), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logFile=nf
		
		if  config.Conf().Mode == "release" {

			logHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			logger = slog.New(logHandler)
		} else {
			writer := io.MultiWriter(os.Stdout, logFile)
			logHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			logger = slog.New(logHandler)
		}
	}
}