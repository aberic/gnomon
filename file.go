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

// Package file 文件操作工具
package gnomon

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type fileCommon struct{}

// PathExists 判断路径是否存在
func (f *fileCommon) PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// ReadFirstLine 从文件中读取第一行并返回字符串数组
func (f *fileCommon) ReadFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	inputString, _ := finReader.ReadString('\n')
	return String().TrimN(inputString), nil
}

// ReadPointLine 从文件中读取指定行并返回字符串数组
func (f *fileCommon) ReadPointLine(filePath string, line int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	lineCount := 1
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			if lineCount == line {
				return inputString, nil
			}
			return "", errors.New("index out of line count")
		}
		if lineCount == line {
			return inputString, nil
		}
		lineCount++
	}
}

// ReadLines 从文件中逐行读取并返回字符串数组
func (f *fileCommon) ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	finReader := bufio.NewReader(file)
	var fileList []string
	for {
		inputString, err := finReader.ReadString('\n')
		//fmt.Println(inputString)
		if err == io.EOF {
			fileList = append(fileList, String().TrimN(inputString))
			break
		}
		fileList = append(fileList, String().TrimN(inputString))
	}
	//fmt.Println("fileList",fileList)
	return fileList, nil
}

// ParentPath 文件父路径
func (f *fileCommon) ParentPath(filePath string) string {
	return filePath[0:strings.LastIndex(filePath, "/")]
}

// CreateAndWrite 创建并写入内容到文件中
// It returns the number of bytes written and an error
func (f *fileCommon) CreateAndWrite(filePath string, data []byte, force bool) (int, error) {
	var (
		file *os.File
		err  error
	)
	exist := f.PathExists(filePath)
	if exist && !force {
		return 0, errors.New("file exist")
	} else if !exist {
		parentPath := f.ParentPath(filePath)
		if err = os.MkdirAll(parentPath, os.ModePerm); nil != err {
			return 0, err
		}
		if file, err = os.Create(filePath); err != nil {
			return 0, err
		}
	} else {
		if file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644); nil != err {
			return 0, err
		}
	}
	defer func() {
		_ = file.Close()
	}()
	// 将数据写入文件中
	//file.WriteString(string(data)) //写入字符串
	if n, err := file.Write(data); nil != err { // 写入byte的slice数据
		return 0, err
	} else {
		return n, nil
	}
}

// LoopDirFromDir 遍历文件夹下的所有子文件夹
func (f *fileCommon) LoopDirFromDir(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Log().Debug("read dir fail", Log().Err(err))
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// LoopAllFileFromDir 遍历文件夹及子文件夹下的所有文件
func (f *fileCommon) LoopAllFileFromDir(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Log().Debug("read dir fail", Log().Err(err))
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = f.LoopAllFileFromDir(fullDir, s)
			if err != nil {
				Log().Debug("read dir fail", Log().Err(err))
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
