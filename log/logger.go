/*
 *  Copyright (c) 2020. aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"encoding/json"
	"fmt"
	"github.com/aberic/gnomon"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// logger 日志工具
type logger struct {
	config *config
	files  map[Level]*filed // files 日志文件输入io对象集合
}

// init log初始化
//
// level 日志级别(debugLevel/infoLevel/warnLevel/ErrorLevel/panicLevel/fatalLevel)
//
// logDir 日志文件目录
//
// maxSize 每个日志文件保存的最大尺寸 单位：M
//
// maxAge 文件最多保存多少天
//
// utc CST & UTC 时间
//
// production 是否生产环境，在生产环境下控制台不会输出任何日志
func (l *logger) init(level Level, logDir string, maxSize, maxAge int, utc bool, production bool) {
	if gnomon.StringIsEmpty(logDir) {
		logDir = gnomon.StringBuild(defaultLogDir)
	}
	if err := os.MkdirAll(logDir, os.ModePerm); nil != err {
		panic(err)
	}
	if nil == l.config {
		l.config = &config{}
	}
	l.config.set(level, logDir, maxSize, maxAge, utc, production)
	if nil == l.files {
		l.files = map[Level]*filed{
			debugLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			infoLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			warnLevel:  {fileIndex: "0", tasks: make(chan string, 1000)},
			errorLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			panicLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			fatalLevel: {fileIndex: "0", tasks: make(chan string, 1000)},
			allLevel:   {fileIndex: "0", tasks: make(chan string, 1000)},
		}
	}
}

// debugSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func (l *logger) debugSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > debugLevel {
		return
	}
	// pc 是uintptr这个返回的是函数指针
	//
	// file 是函数所在文件名目录
	//
	// line 所在行号
	//
	// ok 是否可以获取到信息
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNameDebug, msg, line, debugLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// infoSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
func (l *logger) infoSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > infoLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNameInfo, msg, line, infoLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// warnSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
func (l *logger) warnSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > warnLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNameWarn, msg, line, warnLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// errorSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
func (l *logger) errorSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > errorLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNameError, msg, line, errorLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// panicSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
func (l *logger) panicSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > panicLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNamePanic, msg, line, panicLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// fatalSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
func (l *logger) fatalSkip(skip int, msg string, fields ...FieldInter) {
	if l.config.level > fatalLevel {
		return
	}
	if _, file, line, ok := runtime.Caller(skip); ok {
		l.logStandard(file, logNameFatal, msg, line, fatalLevel, fields...)
	} else {
		panic("log recovery fail")
	}
}

// logStandard 将日志输出到控制台
//
// file 日志触发所在文件地址
//
// levelName 日志级别名称
//
// msg 日志默认输出信息
//
// line 日志触发所在文件的行号
//
// ok 如果无法恢复信息，则为false
//
// level 日志级别
//
// fields 日志输出对象子集
func (l *logger) logStandard(file, levelName, msg string, line int, level Level, fields ...FieldInter) {
	var (
		fileString        string // 即将输出的文件地址信息
		timeString        string // 即将输出的时间信息
		zoneName          string // 即将输出的时区名称
		customContentJSON []byte // 即将输出的自定义日志信息的JSON字节数组
		err               error
	)
	timeNow := time.Now()
	if l.config.utc {
		timeString = timeNow.UTC().Format("2006-01-02 15:04:05.000000")
		zoneName, _ = timeNow.UTC().Zone()
	} else {
		timeString = timeNow.Local().Format("2006-01-02 15:04:05.000000")
		zoneName, _ = timeNow.Local().Zone()
	}
	timeString = strings.Join([]string{timeString, zoneName}, " ")
	logArr := strings.Split(strings.Join([]string{file, strconv.Itoa(line)}, ":"), "/go/src/")
	if len(logArr) > 1 {
		fileString = logArr[1]
	} else {
		fileString = logArr[0]
	}
	customContent := l.customContent(msg, timeString, levelName, fileString, fields...)
	if customContentJSON, err = json.Marshal(customContent); nil != err {
		panic(fmt.Sprint("json Marshal error", err))
	}
	customContentString := string(customContentJSON) // 即将输出的自定义日志信息JSON字符串
	if !l.config.production {
		fmt.Println(timeString, levelName, fileString, customContentString) // 日志输出到控制台
		l.unProduction(customContentString, level)
	} else {
		l.production(customContentString, level)
	}
}

// customContent 用户自定义输出内容集合
//
// msg 日志默认输出信息
//
// timeString 即将输出的时间信息
//
// levelName 日志级别名称
//
// fileString 即将输出的文件地址信息
//
// fields 日志输出对象子集
func (l *logger) customContent(msg, timeString, levelName, fileString string, fields ...FieldInter) map[string]interface{} {
	logCustomContent := make(map[string]interface{}) // 即将输出的日志信息集合
	logCustomContent["msg"] = msg                    // 默认集合第一个参数为msg
	logCustomContent["level"] = strings.ToLower(levelName)
	logCustomContent["time"] = timeString
	logCustomContent["file"] = fileString
	for _, field := range fields {
		if nil == field {
			continue
		}
		logCustomContent[field.GetKey()] = field.GetValue()
	}
	return logCustomContent
}

// unProduction 非生产环境处理策略
//
// customContentString 即将输出的自定义日志信息JSON字符串
//
// level 日志级别
func (l *logger) unProduction(customContentString string, level Level) {
	var stackString string // 即将输出的堆栈信息
	switch level {
	case errorLevel: // 如果是error级别，打印堆栈信息
		stackString = string(debug.Stack())
		fmt.Println(stackString)
	case panicLevel: // 如果是panic级别，打印堆栈信息并执行系统panic方法
		stackString = string(debug.Stack())
		if nil == l.files {
			panic(stackString)
		} else {
			fmt.Println(stackString)
		}
	case fatalLevel: // 如果是fatal级别，打印堆栈信息并终止系统
		stackString = string(debug.Stack())
		fmt.Println(stackString)
		if nil == l.files {
			os.Exit(1)
		}
	}
	if nil == l.files {
		return
	}
	go l.logFile(customContentString, stackString, level)
}

// production 生产环境处理策略
//
// customContentString 即将输出的自定义日志信息JSON字符串
//
// level 日志级别
func (l *logger) production(customContentString string, level Level) {
	if nil == l.files {
		return
	}
	var stackString string
	if level == errorLevel || level == panicLevel || level == fatalLevel {
		stackString = string(debug.Stack())
	}
	go l.logFile(customContentString, stackString, level)
}

// logFile 将日志内容输入文件中存储
//
// customContentString 即将输出的自定义日志信息JSON字符串
//
// stackString 日志堆栈信息
//
// level 日志级别
func (l *logger) logFile(customContentString, stackString string, level Level) {
	var (
		printString string
		err         error
		fd          *filed
	)
	printString = strings.Join([]string{customContentString, stackString}, "\n")
	if fd, err = l.useFiled(level, printString); nil == err {
		fd.tasks <- printString
	}
	if fd, err = l.useFiled(allLevel, printString); nil == err {
		fd.tasks <- printString
	}
}

// useFiled 使用日志文件
//
// level 日志级别
//
// printString 输出字符串
func (l *logger) useFiled(level Level, printString string) (fd *filed, err error) {
	if fd = l.files[level]; fd.file == nil {
		defer fd.lock.Unlock()
		fd.lock.Lock()
		if fd.file == nil {
			var f *os.File
			if f, err = os.OpenFile(l.config.logFilePath(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); nil != err {
				return
			}
			fd.file = f
			if err = l.config.checkFiled(level, fd, int64(len(printString)), false); nil != err {
				return
			}
			go fd.running()
			return
		}
	}
	if err = l.config.checkFiled(level, fd, int64(len(printString)), true); nil != err {
		return
	}
	return
}
