package log

import (
	"github.com/aberic/gnomon"
	"os"
	"strings"
)

const (
	GLogDirEnv         = "LOG_DIR"           // 日志文件目录
	GLogFileMaxSizeEnv = "LOG_FILE_MAX_SIZE" // 每个日志文件保存的最大尺寸 单位：M
	GLogFileMaxAgeEnv  = "LOG_FILE_MAX_AGE"  // 文件最多保存多少天
	GLogUtcEnv         = "LOG_UTC"           // CST & UTC 时间
	GLogLevelEnv       = "LOG_LEVEL"         // 日志级别(debugLevel/infoLevel/warnLevel/ErrorLevel/panicLevel/fatalLevel)
	GLogProductionEnv  = "LOG_PRODUCTION"    // 是否生产环境，在生产环境下控制台不会输出任何日志
)

var (
	logFileDir     string // 日志文件目录
	logFileMaxSize int    // 每个日志文件保存的最大尺寸 单位：M
	logFileMaxAge  int    // 文件最多保存多少天
	logUtc         bool   // CST & UTC 时间
	logLevel       string // 日志级别(debugLevel/infoLevel/warnLevel/ErrorLevel/panicLevel/fatalLevel)
	logProduction  bool   // 是否生产环境，在生产环境下控制台不会输出任何日志
)

func init() {
	logFileDir = gnomon.Env().GetD(GLogDirEnv, os.TempDir())
	logFileMaxSize = gnomon.Env().GetIntD(GLogFileMaxSizeEnv, 1024)
	logFileMaxAge = gnomon.Env().GetIntD(GLogFileMaxAgeEnv, 7)
	logUtc = gnomon.Env().GetBool(GLogUtcEnv)
	logLevel = gnomon.Env().GetD(GLogLevelEnv, "Debug")
	logProduction = gnomon.Env().GetBool(GLogProductionEnv)
	if err := initLog(); nil != err {
		panic(err)
	}
}

func initLog() error {
	if err := gnomon.Log().Init(logFileDir, logFileMaxSize, logFileMaxAge, logUtc); nil != err {
		return err
	}
	var level gnomon.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = gnomon.Log().DebugLevel()
	case "info":
		level = gnomon.Log().InfoLevel()
	case "warn":
		level = gnomon.Log().WarnLevel()
	case "error":
		level = gnomon.Log().ErrorLevel()
	case "panic":
		level = gnomon.Log().PanicLevel()
	case "fatal":
		level = gnomon.Log().FatalLevel()
	default:
		level = gnomon.Log().DebugLevel()
	}
	gnomon.Log().Set(level, logProduction)
	return nil
}

// Param 日志输出子集对象
type Param struct {
	key   string
	value interface{}
}

func (p *Param) GetKey() string {
	return p.key
}

func (p *Param) GetValue() interface{} {
	return p.value
}

// Field 自定义输出KV对象
func Field(key string, value interface{}) *Param {
	return &Param{key: key, value: value}
}

// Err 自定义输出错误
func Err(err error) *Param {
	if nil != err {
		return &Param{key: "error", value: err.Error()}
	}
	return &Param{key: "error", value: nil}
}

// Errs 自定义输出错误
func Errs(msg string) *Param {
	return &Param{key: "error", value: msg}
}

func Debug(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Debug(msg, fields...)
}

func Info(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Info(msg, fields...)
}

func Warn(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Warn(msg, fields...)
}

func Error(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Error(msg, fields...)
}

func Panic(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Panic(msg, fields...)
}

func Fatal(msg string, fields ...gnomon.FieldInter) {
	gnomon.Log().Fatal(msg, fields...)
}
