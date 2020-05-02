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
	"errors"
	"testing"
	"time"
)

var logDir = "tmp/log"

func logDo() {
	Debug("test", Field("1", "2"), Field("2", 3), Field("3", true))
	Debug("test", nil)
	Info("test", Field("1", "2"), Field("2", 3), Field("3", true))
	Warn("test", Field("1", "2"), Field("2", 3), Field("3", true))
	Error("test", Field("1", "2"), Field("2", 3), Field("3", true), Err(errors.New("yes")))
	DebugSkip(1, "test", Field("1", "2"), Field("2", 3), Field("3", true))
	DebugSkip(1, "test", nil)
	InfoSkip(1, "test", Field("1", "2"), Field("2", 3), Field("3", true))
	WarnSkip(1, "test", Field("1", "2"), Field("2", 3), Field("3", true))
	ErrorSkip(1, "test", Field("1", "2"), Field("2", 3), Field("3", true), Err(errors.New("yes")))
}

func TestLogCommon_Debug(t *testing.T) {
	Set(DebugLevel(), logDir, 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}

func TestLogCommon_Info(t *testing.T) {
	Set(InfoLevel(), logDir, 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}

func TestLogCommon_Warn(t *testing.T) {
	Set(WarnLevel(), logDir, 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}

func TestLogCommon_Error(t *testing.T) {
	Set(ErrorLevel(), logDir, 1, 1, false, false)
	logDo()
}

func TestLogCommon_Panic(t *testing.T) {
	Set(PanicLevel(), logDir, 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}

func TestLogCommon_Fatal(t *testing.T) {
	Set(FatalLevel(), logDir, 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}

func TestLogCommon_Fatal_BigStorage(t *testing.T) {
	Set(DebugLevel(), logDir, 1, 1, false, true)
	for i := 0; i < 100000; i++ {
		go logDo()
	}
	time.Sleep(10 * time.Second)
}
