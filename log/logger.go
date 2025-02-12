package log

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/config"
	"github.com/mitchellh/go-homedir"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// global sys var
var (
	GlobalLogger *zap.Logger
	Level        zap.AtomicLevel = zap.NewAtomicLevel()
	logCnf       *config.LogConfig
)

var (
	// after the app close, we should sync the message to io.Writer
	Sync func() error
)

var (
	Info  func(msg string, fields ...zap.Field)
	Debug func(msg string, fields ...zap.Field)
	Error func(msg string, fields ...zap.Field)
	Warn  func(msg string, fields ...zap.Field)
	Fatal func(msg string, fields ...zap.Field)

	// work like stand pkg sys.Printf
	Infof  func(template string, args ...interface{})
	Warnf  func(template string, args ...interface{})
	Debugf func(template string, args ...interface{})
	Errorf func(template string, args ...interface{})
	Fataf  func(template string, args ...interface{})
)

func newStdOutLogCore() zapcore.Core {

	// std use the console encoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		Level,
	)
}

func newLogFileLogCore() zapcore.Core {
	// Set the fd
	if len(logCnf.LogDir) == 0 {
		return zapcore.NewNopCore()
	}
	path, err := homedir.Expand(logCnf.LogDir)
	if err != nil {
		return zapcore.NewNopCore()
	}
	_, err = os.Stat(path)
	notexist := os.IsNotExist(err)
	if notexist {
		fmt.Println("Initializing log path at ", logCnf.LogDir)
		err = os.MkdirAll(logCnf.LogDir, 0755)
		fmt.Println("Initializing log path err ", err)
		if err != nil && !os.IsExist(err) {
			return zapcore.NewNopCore()
		}
	}

	logFileName := logCnf.LogDir + "/rock-cloud.log"

	logFileFD := lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    logCnf.MaxSize,
		MaxAge:     logCnf.MaxAge,
		MaxBackups: logCnf.MaxBackups,
		LocalTime:  logCnf.LocalTime,
	}
	Level.SetLevel(zapcore.Level(logCnf.Level))
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
	}

	// sys-file the console encoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(logCfg.EncoderConfig),
		zapcore.AddSync(&logFileFD),
		Level,
	)
}

func initLogger() {
	trace := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})

	GlobalLogger = zap.New(zapcore.NewTee(
		newStdOutLogCore(),
		newLogFileLogCore(),
	),
		zap.AddCaller(),
		zap.AddStacktrace(trace),
	)

	Info = GlobalLogger.Info
	Debug = GlobalLogger.Debug
	Warn = GlobalLogger.Warn
	Error = GlobalLogger.Error
	Fatal = GlobalLogger.Fatal

	// work like stand pkg sys.Printf
	Infof = GlobalLogger.Sugar().Infof
	Warnf = GlobalLogger.Sugar().Warnf
	Debugf = GlobalLogger.Sugar().Debugf
	Errorf = GlobalLogger.Sugar().Errorf
	Fataf = GlobalLogger.Sugar().Fatalf

	// the sync method
	Sync = GlobalLogger.Sync
}

func ConfigLog(cnf *config.LogConfig) {
	logCnf = cnf
	fmt.Println("rockx MaxAge", logCnf.MaxAge, "MaxSize", logCnf.MaxSize, "LogDir", logCnf.LogDir)
	initLogger()
}
