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
	"os"
	"path/filepath"
	"sync"
)

var (
	logPhysical   *logger
	defaultLogDir = filepath.Join(os.TempDir(), "log")
	once          sync.Once // once log对象只会被初始化一次
)

func init() {
	logPhysical = &logger{}
	logPhysical.init(debugLevel, defaultLogDir, 1024, 7, false, false)
}

// Set log初始化
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
func Set(level Level, logDir string, maxSize, maxAge int, utc bool, production bool) {
	once.Do(func() {
		logPhysical.init(level, logDir, maxSize, maxAge, utc, production)
	})
}

// Debug 输出指定级别日志
func Debug(msg string, fields ...FieldInter) {
	logPhysical.debugSkip(2, msg, fields...)
}

// Info 输出指定级别日志
func Info(msg string, fields ...FieldInter) {
	logPhysical.infoSkip(2, msg, fields...)
}

// Warn 输出指定级别日志
func Warn(msg string, fields ...FieldInter) {
	logPhysical.warnSkip(2, msg, fields...)
}

// Error 输出指定级别日志
func Error(msg string, fields ...FieldInter) {
	logPhysical.errorSkip(2, msg, fields...)
}

// Panic 输出指定级别日志
func Panic(msg string, fields ...FieldInter) {
	logPhysical.panicSkip(2, msg, fields...)
}

// Fatal 输出指定级别日志
func Fatal(msg string, fields ...FieldInter) {
	logPhysical.fatalSkip(2, msg, fields...)
}

// DebugSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func DebugSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.debugSkip(skip, msg, fields...)
}

// InfoSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func InfoSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.infoSkip(skip, msg, fields...)
}

// WarnSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func WarnSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.warnSkip(skip, msg, fields...)
}

// ErrorSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func ErrorSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.errorSkip(skip, msg, fields...)
}

// PanicSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func PanicSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.panicSkip(skip, msg, fields...)
}

// FatalSkip 输出指定级别日志
//
// skip 提升的堆栈帧数，0-当前函数，1-上一层函数。如果经封装调用该方法，默认2，否则默认1
//
// msg 日志默认输出信息
//
// fields 日志输出对象子集
func FatalSkip(skip int, msg string, fields ...FieldInter) {
	logPhysical.fatalSkip(skip, msg, fields...)
}

// Field 自定义输出KV对象
func Field(key string, value interface{}) FieldInter {
	return &field{key: key, value: value}
}

// Err 自定义输出错误
func Err(err error) FieldInter {
	if nil != err {
		return &field{key: "error", value: err.Error()}
	}
	return &field{key: "error", value: nil}
}

// Errs 自定义输出错误
func Errs(msg string) FieldInter {
	return &field{key: "error", value: msg}
}

// DebugLevel logs are typically voluminous, and are usually disabled in production.
func DebugLevel() Level {
	return debugLevel
}

// InfoLevel is the default logging priority.
func InfoLevel() Level {
	return infoLevel
}

// WarnLevel logs are more important than Info, but don't need individual human review.
func WarnLevel() Level {
	return warnLevel
}

// ErrorLevel logs are high-priority. If an application is running smoothly,
// it shouldn't generate any error-level logs.
func ErrorLevel() Level {
	return errorLevel
}

// PanicLevel logs a message, then panics.
func PanicLevel() Level {
	return panicLevel
}

// FatalLevel logs a message, then panics.
func FatalLevel() Level {
	return fatalLevel
}
