package bootstrap

import (
	"fmt"
	"logservice/global"
	"logservice/utils"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
初始化本地服务日志器
*/
func InitializeLocalLogger() {
	lS := global.App.Config.LocalServiceLog
	bS := global.App.Config.BaseServiceLog
	core := genCore(lS.JsonFormat, lS.Level, lS.RootDir, lS.Info, lS.Error, bS.MaxSize, bS.MaxAge, bS.MaxBackups)
	logger := zap.New(core, zap.AddCaller()).Named(lS.ServiceName)
	global.App.LogsServiceLogger = *logger
}

/*
初始化用户微服务日志器
*/
func InitializeUserServiceLogger() {
	uS := global.App.Config.UserServiceLog
	bS := global.App.Config.BaseServiceLog
	core := genCore(uS.JsonFormat, uS.Level, uS.RootDir, uS.Info, uS.Error, bS.MaxSize, bS.MaxAge, bS.MaxBackups)
	logger := zap.New(core, zap.AddCaller()).Named(uS.ServiceName)
	global.App.UserServiceLogger = *logger
}

/*
定制zapCore方法
zapCore 核心三组件 encoder LevelEnabler WriteSyncer
请注意，`DPanic`、`Panic`和`Fatal`级别的日志记录在功能上有所不同：
- `DPanic`会在开发模式下`panic`，但在生产模式下不会。
- `Panic`级别的日志会记录日志后调用`panic`。
- `Fatal`级别的日志会记录日志后立即调用`os.Exit(1)`终止程序。
*/
func genCore(isJson bool, levelStr string, rootdir string, infon string, errn string, maxsize int, maxage int, maxbackups int) zapcore.Core {
	encoder := getEncoder(isJson, global.App.Config.ServiceInfo.Env)
	// 定义级别,baseErrorLevel 默认错误等级,InfoLevel 表示日志等级
	errorLevel := getLevel(global.App.Config.BaseServiceLog.DefaultErrorL)
	infoLevel := getLevel(levelStr)
	if errorLevel < infoLevel {
		panic("错误日志等级低于默认日志等级，错误日志不会被记录")
	}
	consoleWS := zapcore.AddSync(os.Stdout)
	info_file_name := global.App.Config.BaseServiceLog.RootDir + rootdir + infon
	error_file_name := global.App.Config.BaseServiceLog.RootDir + rootdir + errn
	infofileWS := getFileWriteSyncer(info_file_name, maxsize, maxage, maxbackups)
	errfileWS := getFileWriteSyncer(error_file_name, maxsize, maxage, maxbackups)

	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infofileWS, consoleWS), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < errorLevel && l >= infoLevel
	}))
	errCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errfileWS, consoleWS), errorLevel)

	core := zapcore.NewTee(infoCore, errCore)
	return core
}

/*
根据参数不同，生成不同的encoder
这里参数主要是	NameKey: "logger" 进行更改，方便在控制台中区分不同服务的日志信息
参数解释：
serverName： 区分不同服务日志
*/
func getEncoder(isJson bool, env string) (encoder zapcore.Encoder) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Now().Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(env + "." + l.String())
	}
	if isJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

/*
日志级别相关配置
不同服务需要的日志级别可能不同，假设有两个服务，一个是记录用户签到记录用户活跃度的服务，另外一个是支付服务，对于前者，可能只需要输出日志等级为error的信息，后者可能需要输出info以及以上的日志信息
这里函数将字符串转为zap 的 level 类型
*/
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
	fmt.Println("文件名",filename)
	fmt.Println("文件路径",filepath.Dir(filename))
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
