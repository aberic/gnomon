/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"os"
	"strconv"
	"strings"
)

// EnvGet 获取环境变量 envName 的值
//
// envName 环境变量名称
func EnvGet(envName string) string {
	return os.Getenv(envName)
}

// EnvGetD 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func EnvGetD(envName string, defaultValue string) string {
	env := EnvGet(envName)
	if StringIsEmpty(env) {
		return defaultValue
	}
	return env
}

// EnvGetInt 获取环境变量 envName 的值
//
// envName 环境变量名称
func EnvGetInt(envName string) (int, error) {
	return strconv.Atoi(os.Getenv(envName))
}

// EnvGetIntD 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func EnvGetIntD(envName string, defaultValue int) int {
	if i, err := strconv.Atoi(os.Getenv(envName)); nil == err {
		return i
	}
	return defaultValue
}

// EnvGetInt64 获取环境变量 envName 的值
//
// envName 环境变量名称
func EnvGetInt64(envName string) (int64, error) {
	return strconv.ParseInt(EnvGet(envName), 10, 64)
}

// EnvGetInt64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func EnvGetInt64D(envName string, defaultValue int64) int64 {
	if i, err := strconv.ParseInt(EnvGet(envName), 10, 64); nil == err {
		return i
	}
	return defaultValue
}

// EnvGetUint64 获取环境变量 envName 的值
//
// envName 环境变量名称
func EnvGetUint64(envName string) (uint64, error) {
	return strconv.ParseUint(EnvGet(envName), 10, 64)
}

// EnvGetUint64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func EnvGetUint64D(envName string, defaultValue uint64) uint64 {
	if i, err := strconv.ParseUint(EnvGet(envName), 10, 64); nil == err {
		return i
	}
	return defaultValue
}

// EnvGetFloat64 获取环境变量 envName 的值
//
// envName 环境变量名称
func EnvGetFloat64(envName string) (float64, error) {
	return strconv.ParseFloat(EnvGet(envName), 64)
}

// EnvGetFloat64D 获取环境变量 envName 的值
//
// envName 环境变量名称
//
// defaultValue 环境变量为空时的默认值
func EnvGetFloat64D(envName string, defaultValue float64) float64 {
	if i, err := strconv.ParseFloat(EnvGet(envName), 64); nil == err {
		return i
	}
	return defaultValue
}

// EnvGetBool 获取环境变量 envName 的 bool 值
//
// envName 环境变量名称
func EnvGetBool(envName string) bool {
	return strings.EqualFold(os.Getenv(envName), "true")
}
