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
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/robfig/cron"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// config 日志工具
type config struct {
	logDir      string     // logDir 日志文件目录
	maxSizeByte int64      // maxSizeByte 每个日志文件保存的最大尺寸 单位：byte
	maxAge      int        // maxAge 文件最多保存多少天
	level       Level      // level 日志级别
	production  bool       // 生产环境，该模式下控制台不会输出任何日志
	utc         bool       // CST & UTC 时间
	date        string     // date 当前日志文件后缀日期
	job         *cron.Cron // job 日志定时清理任务
}

// set config配置设置
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
func (l *config) set(level Level, logDir string, maxSize, maxAge int, utc bool, production bool) {
	l.logDir = logDir
	l.utc = utc
	if maxSize < 1 {
		maxSize = 1
	}
	l.maxSizeByte = int64(maxSize * 1024 * 1024)
	if maxAge < 1 {
		maxAge = 1
	}
	l.maxAge = maxAge
	l.level = debugLevel
	l.production = false
	if utc {
		l.date = time.Now().UTC().Format("20060102")
	} else {
		l.date = time.Now().Local().Format("20060102")
	}
	switch level {
	default:
		l.level = debugLevel
	case debugLevel, infoLevel, warnLevel, errorLevel, panicLevel, fatalLevel:
		l.level = level
	}
	l.production = production
	// todo reset
	l.job = cron.New()
	go l.checkMaxAge()
}

// checkMaxAge 遍历并检查文件是否达到保存天数，达到则删除
func (l *config) checkMaxAge() {
	// 每隔5秒执行一次：*/5 * * * * ?
	// 每隔1分钟执行一次：0 */1 * * * ?
	// 每天23点执行一次：0 0 23 * * ?
	// 每天凌晨1点执行一次：0 0 1 * * ?
	// 每月1号凌晨1点执行一次：0 0 1 1 * ?
	// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
	// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	err := l.job.AddFunc(strings.Join([]string{"0 0 0 */", strconv.Itoa(l.maxAge), " * ?"}, ""), func() {
		var timeDate string
		if l.utc {
			timeDate = time.Now().UTC().Format("20060102")
		} else {
			timeDate = time.Now().Local().Format("20060102")
		}
		logDirs, _ := gnomon.FileLoopDirs(l.logDir)
		for _, dirName := range logDirs {
			if strings.Contains(dirName, timeDate) {
				_ = os.RemoveAll(dirName)
			}
		}
	})
	if nil != err {
		time.Sleep(time.Second)
		l.checkMaxAge()
	} else {
		l.job.Start()
	}
}

// logFilePath 日志文件路径
func (l *config) logFilePath(fd *filed, level Level) string {
	parentPath := filepath.Join(l.logDir, l.date)
	if exist := gnomon.FilePathExists(parentPath); !exist {
		if err := os.MkdirAll(parentPath, os.ModePerm); nil != err {
			panic(fmt.Sprint("path mkdirAll error", err))
		}
	}
	return filepath.Join(parentPath, l.levelFileName(fd, level))
}

// levelFileName 包含日志类型的日志文件名称
func (l *config) levelFileName(fd *filed, level Level) string {
	switch level {
	case debugLevel:
		return l.logFileName("debug_", fd.fileIndex)
	case infoLevel:
		return l.logFileName("info_", fd.fileIndex)
	case warnLevel:
		return l.logFileName("warn_", fd.fileIndex)
	case errorLevel:
		return l.logFileName("error_", fd.fileIndex)
	case panicLevel:
		return l.logFileName("panic_", fd.fileIndex)
	case fatalLevel:
		return l.logFileName("fatal_", fd.fileIndex)
	}
	return l.logFileName("log_", fd.fileIndex)
}

// logFileName 不包含日志类型的日志文件名称
func (l *config) logFileName(name, index string) string {
	return gnomon.StringBuild(name, l.date, "-", index, ".log")
}

// checkFiled 检查日志文件是否可用
//
// 如果当前正在使用的日志文件已经达到单个文件大小上限，则通过后缀++的方式将内容写入新的文件中
//
// level 日志级别
//
// fd 日志文件操作对象
//
// printStringLength 输出到文件中字节数长度
//
// lock 该操作是否需要给filed文件对象上锁。如果是复用对象，则需要上锁；如果是新建对象，则新建过程中本身就已经上锁，此处无需锁定
func (l *config) checkFiled(level Level, fd *filed, printStringLength int64) (err error) {
	var ret int64
	if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err { // 当前文件已用字节数
		return
	}
	// 文件字节数上限 - 当前文件已用字节数 - 即将消耗文件字节数
	if l.maxSizeByte-ret-printStringLength < 0 { // 如果小于0，则说明文件长度不足，需要新建文件
		defer fd.lock.Unlock()
		fd.lock.Lock()
		return l.findAvailableFile(level, fd, printStringLength)
	}
	return
}

// findAvailableFile 查找可用日志文件
//
// 如果当前正在使用的日志文件已经达到单个文件大小上限，则通过后缀++的方式将内容写入新的文件中
//
// level 日志级别
//
// fd 日志文件操作对象
func (l *config) findAvailableFile(level Level, fd *filed, printStringLength int64) (err error) {
	var ret int64
	if ret, err = fd.file.Seek(0, io.SeekEnd); nil != err {
		return
	}
	// 文件字节数上限 - 当前文件已用字节数 - 即将消耗文件字节数
	if l.maxSizeByte-ret-printStringLength < 0 { // 如果小于0，则说明文件长度不足，需要新建文件
		index, _ := strconv.Atoi(fd.fileIndex) // 获取旧文件的序列号
		fd.fileIndex = strconv.Itoa(index + 1) // 赋值新文件的序列号
		// 创建并赋值新文件
		if fd.file, err = os.OpenFile(l.logFilePath(fd, level), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); nil != err {
			return
		}
	}
	return nil
}
