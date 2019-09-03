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

// Package env 环境变量工具
package gnomon

import (
	"os"
	"strings"
)

type envCommon struct{}

// GetEnv 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) GetEnv(envName string) string {
	return os.Getenv(envName)
}

// GetEnvBool 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
func (e *envCommon) GetEnvBool(envName string) bool {
	return strings.EqualFold(os.Getenv(envName), "true")
}

// GetEnvDefault 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetEnvDefault(envName string, defaultValue string) string {
	env := e.GetEnv(envName)
	if String().IsEmpty(env) {
		return defaultValue
	}
	return env
}

// GetEnvBoolDefault 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetEnvBoolDefault(envName string, defaultValue bool) bool {
	env := e.GetEnv(envName)
	if String().IsEmpty(env) {
		return defaultValue
	}
	return strings.EqualFold(os.Getenv(envName), "true")
}
