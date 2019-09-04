package gnomon

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
	AllLevel
)

const (
	logNameDebug = "DEBUG"
	logNameInfo  = "INFO "
	logNameWarn  = "WARN "
	logNameError = "ERROR"
	logNamePanic = "PANIC"
	logNameFatal = "FATAL"
)

type logCommon struct {
	logDir           string           // logDir 日志文件目录
	maxSizeByte      int64            // maxSizeByte 每个日志文件保存的最大尺寸 单位：byte
	maxBackups       int              // maxBackups 日志文件最多保存多少个备份
	maxAge           int              // maxAge 文件最多保存多少天
	compress         bool             // compress 是否压缩
	files            map[Level]*filed // files 日志文件输入io对象集合
	level            Level            // level 日志级别
	production       bool             // 生产环境，该模式下控制台不会输出任何日志
	utc              bool             // CST & UTC 时间
	date             string           // date 当前日志文件后缀日期
	mkRootDirSuccess bool             // mkRootDirSuccess 是否成功初始化log对象
	once             sync.Once        // once log对象只会被初始化一次
}

// Init log初始化
//
// logDir 日志文件目录
//
// maxSize 每个日志文件保存的最大尺寸 单位：M
//
// maxBackups 日志文件最多保存多少个备份
//
// maxAge 文件最多保存多少天
//
// compress 是否压缩
//
// utc CST & UTC 时间
func (l *logCommon) Init(logDir string, maxSize, maxBackups, maxAge int, utc, compress bool) {
	l.once.Do(func() {
		if err := os.MkdirAll(logDir, os.ModePerm); nil != err {
			log.Println(err)
			return
		}
		l.mkRootDirSuccess = true
		l.logDir = logDir
		l.utc = utc
		l.maxSizeByte = int64(maxSize * 1024 * 1024)
		l.maxBackups = maxBackups
		l.maxAge = maxAge
		l.compress = compress
		l.files = map[Level]*filed{
			DebugLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			InfoLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			WarnLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			ErrorLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			PanicLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			FatalLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			AllLevel:   {fileIndex: "0", tasks: make(chan string, 1000)},
		}
		if utc {
			l.date = time.Now().UTC().Format("20060102")
		} else {
			l.date = time.Now().Local().Format("20060102")
		}
	})
}

func (l *logCommon) SetLevel(level Level, production bool) {
	l.level = level
	l.production = production
}

func (l *logCommon) Debug(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameDebug, msg, line, ok, DebugLevel, fields...)
	}
}

func (l *logCommon) Info(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameInfo, msg, line, ok, InfoLevel, fields...)
	}
}

func (l *logCommon) Warn(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameWarn, msg, line, ok, WarnLevel, fields...)
	}
}

func (l *logCommon) Error(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameError, msg, line, ok, ErrorLevel, fields...)
	}
}

func (l *logCommon) Panic(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNamePanic, msg, line, ok, PanicLevel, fields...)
	}
}

func (l *logCommon) Fatal(msg string, fields ...*field) {
	if _, file, line, ok := runtime.Caller(1); ok {
		l.logStandard(file, logNameFatal, msg, line, ok, FatalLevel, fields...)
	}
}

// logStandard 输出日志到控制台
func (l *logCommon) logStandard(file, levelName, msg string, line int, ok bool, level Level, fields ...*field) {
	if l.level > level {
		return
	}
	var (
		fileString  string
		timeString  string
		stackString string
	)
	timeNow := time.Now()
	if l.utc {
		timeString = timeNow.UTC().String()
	} else {
		timeString = timeNow.Local().String()
	}
	fileString = strings.Join([]string{file, strconv.Itoa(line)}, ":")
	if !l.production {
		var (
			commandJson []byte
			err         error
		)
		logCommand := make(map[string]interface{})
		logCommand["msg"] = msg
		for _, field := range fields {
			logCommand[field.key] = field.value
		}
		if commandJson, err = json.Marshal(logCommand); nil != err {
			log.Println("json Marshal error")
			return
		}
		fmt.Println(timeString, levelName, fileString, string(commandJson))
		switch levelName {
		case logNameError, logNamePanic, logNameFatal:
			stackString = string(debug.Stack())
			fmt.Println(stackString)
		}
	}
	_ = pool().submitField(func(timeString, fileString, stackString, levelName, msg string, level Level, fields ...*field) {
		l.logFile(timeString, fileString, stackString, levelName, msg, level, fields...)
	}, timeString, fileString, stackString, levelName, msg, level, fields...)
}

func (l *logCommon) logFile(timeString, fileString, stackString, levelName, msg string, level Level, fields ...*field) {
	var (
		mapJson     []byte
		printString string
		err         error
		fd          *filed
	)
	logMap := make(map[string]interface{})
	logMap["level"] = strings.ToLower(levelName)
	logMap["time"] = timeString
	logMap["file"] = fileString
	logMap["msg"] = msg
	for _, field := range fields {
		logMap[field.key] = field.value
	}
	if mapJson, err = json.Marshal(logMap); nil != err {
		log.Println("json Marshal error")
		return
	}
	switch levelName {
	case logNameError, logNamePanic, logNameFatal:
		if String().IsEmpty(stackString) {
			stackString = string(debug.Stack())
		}
		printString = strings.Join([]string{string(mapJson), stackString}, "\n")
	default:
		printString = strings.Join([]string{string(mapJson), "\n"}, "")
	}
	if fd, err = l.useFiled(level, printString); nil == err {
		fd.tasks <- printString
	}
	if fd, err = l.useFiled(AllLevel, printString); nil == err {
		fd.tasks <- printString
	}
}

func (l *logCommon) useFiled(level Level, printString string) (fd *filed, err error) {
	if fd = l.files[level]; fd.file == nil {
		defer fd.lock.Unlock()
		fd.lock.Lock()
		if fd.file == nil {
			var f *os.File
			if f, err = os.OpenFile(l.path(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); nil != err {
				return
			}
			fd.file = f
			if err = l.checkFiled(level, fd, int64(len(printString)), false); nil != err {
				return
			}
			err = pool().submit(func() {
				fd.running()
			})
			return
		}
	}
	if err = l.checkFiled(level, fd, int64(len(printString)), true); nil != err {
		return
	}
	return
}

func (l *logCommon) checkFiled(level Level, fd *filed, printStringLength int64, lock bool) (err error) {
	var ret int64
	if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
		return
	}
	if l.maxSizeByte-ret-printStringLength < 0 {
		if lock {
			defer fd.lock.Unlock()
			fd.lock.Lock()
			if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
				return
			}
			if l.maxSizeByte-ret-printStringLength < 0 {
				index, _ := strconv.Atoi(fd.fileIndex)
				fd.fileIndex = strconv.Itoa(index + 1)
				fd.file, err = os.OpenFile(l.path(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			}
		} else {
			if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
				return
			}
			if l.maxSizeByte-ret-printStringLength < 0 {
				index, _ := strconv.Atoi(fd.fileIndex)
				fd.fileIndex = strconv.Itoa(index + 1)
				fd.file, err = os.OpenFile(l.path(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			}
		}
	}
	return
}

func (l *logCommon) path(fd *filed, level Level) string {
	return filepath.Join(l.logDir, l.levelFileName(fd, level))
}

func (l *logCommon) levelFileName(fd *filed, level Level) string {
	switch level {
	case DebugLevel:
		return l.logFileName("debug_", fd.fileIndex)
	case InfoLevel:
		return l.logFileName("info_", fd.fileIndex)
	case WarnLevel:
		return l.logFileName("warn_", fd.fileIndex)
	case ErrorLevel:
		return l.logFileName("error_", fd.fileIndex)
	case PanicLevel:
		return l.logFileName("panic_", fd.fileIndex)
	case FatalLevel:
		return l.logFileName("fatal_", fd.fileIndex)
	}
	return l.logFileName("log_", fd.fileIndex)
}

func (l *logCommon) logFileName(name, index string) string {
	return strings.Join([]string{name, l.date, "-", index, ".log"}, "")
}

type filed struct {
	fileIndex string // fileIndex 日志文件相同日期编号，根据文件新建规则确定
	file      *os.File
	tasks     chan string
	check     chan int8
	lock      sync.Mutex // lock 每次做io开销的安全锁
}

func (f *filed) running() {
	to := time.NewTimer(60 * time.Second)
	for {
		select {
		case task := <-f.tasks:
			to.Reset(time.Second)
			if _, err := f.file.WriteString(task); nil != err {
				panic(err)
			}
		case <-f.check: // 文件存储满，切换文件写入

		case <-to.C:
			_ = f.file.Close()
			f.file = nil
			return
		}
	}
}

type field struct {
	key   string
	value interface{}
}
