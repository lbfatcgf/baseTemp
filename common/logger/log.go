package logger

import (
	"fmt"

	"codeup.aliyun.com/67c7c688484ca2f0a13acc04/baseTemp/common/config"
	"codeup.aliyun.com/67c7c688484ca2f0a13acc04/baseTemp/tools"

	"log/slog"
	"os"
	fp "path/filepath"
	"time"
)

var (
	logger     *slog.Logger
	cmdLogger  *slog.Logger
	logMsgChan chan *logMsg
	logFile    *os.File
	logdir     string
)

type logMsg struct {
	Level   slog.Level
	Message string
	arg     []any
}

func LogInfo(msg string, arg ...any) {
	if cmdLogger != nil {
		cmdLogger.Info(msg, arg...)
	}
	logMsgChan <- &logMsg{Level: slog.LevelInfo, Message: msg, arg: arg}
}

func LogError(msg string, arg ...any) {
	if cmdLogger != nil {
		cmdLogger.Error(msg, arg...)
	}
	logMsgChan <- &logMsg{Level: slog.LevelError, Message: msg, arg: arg}
}

func LogWarn(msg string, arg ...any) {
	if cmdLogger != nil {
		cmdLogger.Warn(msg, arg...)
	}
	logMsgChan <- &logMsg{Level: slog.LevelWarn, Message: msg, arg: arg}
}
func LogDebug(msg string, arg ...any) {
	if cmdLogger != nil {
		cmdLogger.Debug(msg, arg...)
	}
	logMsgChan <- &logMsg{Level: slog.LevelDebug, Message: msg, arg: arg}
}
func InitLog(mode string) {

	logMsgChan = make(chan *logMsg, 1000)
	logdir := config.Conf().LogDir
	if len(logdir) == 0 {
		exepath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		logdir = fp.Join(fp.Dir(exepath), tools.SafeFilePath("logs"))
	}
	// fmt.Println(fp.Join(logdir,Conf().Name+"-"+time.Now().Format("2006-01-02")+".log"))
	if err := os.MkdirAll(logdir, 0755); err != nil {
		panic(err)
	}
	if config.Conf().Mode != "release" {
		fmt.Println("log dir:", logdir)
	}
	openFile, err := os.OpenFile(fp.Join(logdir, tools.SafeFilePath(config.Conf().Name+"-"+time.Now().Format("2006-01-02")+".log")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)

	}
	logFile = openFile
	logHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger = slog.New(logHandler)
	if mode != "release" {
		cmdLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	startLog()

}

func startLog() {
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
	// logdir := config.Conf().LogDir
	newFilePtah := fp.Join(logdir, tools.SafeFilePath(config.Conf().Name+"-"+time.Now().Format("2006-01-02")+".log"))
	if newFilePtah != logFile.Name() {
		logFile.Close()
		nf, err := os.OpenFile(fp.Join(config.Conf().LogDir, newFilePtah), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logFile = nf

		logHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(logHandler)

	}
}
