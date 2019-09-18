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
	"strconv"
	"strings"
)

type envCommon struct{}

// Get 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) Get(envName string) string {
	return os.Getenv(envName)
}

// GetD 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetD(envName string, defaultValue string) string {
	env := e.Get(envName)
	if String().IsEmpty(env) {
		return defaultValue
	}
	return env
}

// GetInt 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) GetInt(envName string) (int, error) {
	return strconv.Atoi(os.Getenv(envName))
}

// GetIntD 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetIntD(envName string, defaultValue int) int {
	if i, err := strconv.Atoi(os.Getenv(envName)); nil == err {
		return i
	}
	return defaultValue
}

// GetInt64 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) GetInt64(envName string) (int64, error) {
	return strconv.ParseInt(e.Get(envName), 10, 64)
}

// GetInt64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetInt64D(envName string, defaultValue int64) int64 {
	if i, err := strconv.ParseInt(e.Get(envName), 10, 64); nil == err {
		return i
	}
	return defaultValue
}

// GetUint64 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) GetUint64(envName string) (uint64, error) {
	return strconv.ParseUint(e.Get(envName), 10, 64)
}

// GetUint64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetUint64D(envName string, defaultValue uint64) uint64 {
	if i, err := strconv.ParseUint(e.Get(envName), 10, 64); nil == err {
		return i
	}
	return defaultValue
}

// GetFloat64 获取环境变量 envName 的值
//
// envName 环境变量名称
func (e *envCommon) GetFloat64(envName string) (float64, error) {
	return strconv.ParseFloat(e.Get(envName), 64)
}

// GetFloat64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func (e *envCommon) GetFloat64D(envName string, defaultValue float64) float64 {
	if i, err := strconv.ParseFloat(e.Get(envName), 64); nil == err {
		return i
	}
	return defaultValue
}

// GetBool 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
func (e *envCommon) GetBool(envName string) bool {
	return strings.EqualFold(os.Getenv(envName), "true")
}
