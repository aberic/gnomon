/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
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
 *
 */

package log

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logImpl := Logger{
		Config: &Config{
			Level:      DebugLevel,
			MaxSize:    128,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
		}}

	testLogger := logImpl.NewCustom("./logs/ho.log", DebugLevel, 128, 30, 30, true, "ho")
	testLogger.Info("log 初始化成功1")

	logImpl.Conf(&Config{
		Level:      DebugLevel,
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	})

	testLogger = logImpl.New("./logs/ho.log", "ho")
	testLogger.Info("log 初始化成功2")

	testLogger = GetLogInstance().New("./logs/instance.log", "instance")
	testLogger.Info("log 初始化成功3")

	testLogger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
}
