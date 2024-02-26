package bootstrap

import (
	"auth/global"
	"auth/utils"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
	这里封装服务自带的log，在连接rabbitMQ 或者 发送log 失败时，做记录
*/

func InitLocalLogger() {
	localLog := global.App.Config.LocalLogs
	core := genCore(localLog.IsJson, localLog.DefaultLevel, localLog.Dir, localLog.Logfilename, localLog.Max_size, localLog.Max_size, localLog.Max_backups)
	logger := zap.New(core, zap.AddCaller())
	global.App.LocalLogger = logger
}

func genCore(isJson bool, levelStr string, rootdir string, infon string, maxsize int, maxage int, maxbackups int) zapcore.Core {
	encoder := getEncoder(isJson)
	infoLevel := getLevel(levelStr)
	consoleWS := zapcore.AddSync(os.Stdout)
	info_file_name := rootdir + infon
	infofileWS := getFileWriteSyncer(info_file_name, maxsize, maxage, maxbackups)
	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infofileWS, consoleWS), infoLevel)
	return infoCore
}

func getEncoder(isJson bool) (encoder zapcore.Encoder) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Now().Format("2006-01-02 15:04:05.000"))
	}
	if isJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

// 日志级别相关配置
func getLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

/*
	自定义写入器
	一般有两种，写入文件或者控制台
*/

func getFileWriteSyncer(filename string, maxsize int, maxage int, maxbackups int) zapcore.WriteSyncer {
	utils.CreateDir(filepath.Dir(filename))
	file := &lumberjack.Logger{
		Filename:   filename,   // 文件名称和路径
		MaxSize:    maxsize,    // 日志文件的最大大小（以MB为单位）
		MaxAge:     maxage,     // 旧日志文件保留的最大天数
		MaxBackups: maxbackups, // 保留的最大旧日志文件数量
		Compress:   true,       // 对旧的日志文件进行压缩
	}
	return zapcore.AddSync(file)
}
