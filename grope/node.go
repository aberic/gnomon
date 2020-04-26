/*
 * Copyright (c) 2020. Aberic - All Rights Reserved.
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

package grope

import (
	"fmt"
	"github.com/aberic/gnomon"
	"strings"
	"sync"
)

func newNode(filters ...Filter) *node {
	return &node{
		root:      true,
		filters:   filters,
		nextNodes: []*node{},
	}
}

func nextNode(n *node, patternPiece string) *node {
	return &node{
		root:         false,
		patternPiece: patternPiece,
		filters:      n.filters,
		preNode:      n,
		nextNodes:    []*node{},
	}
}

type node struct {
	root         bool     // 是否根结点
	pattern      string   // /a/b/:c/d/:e/:f/g
	patternPiece string   // a || ?
	method       string   // eg:http.MethodGet
	handler      Handler  // 待实现接收请求方法
	filters      []Filter // 过滤器/拦截器数组
	preNode      *node
	nextNodes    []*node
	lockNode     sync.Mutex
}

// add
//
// pattern /a/b/:c/d/:e/:f/g
//
// method eg:http.MethodGet
func (n *node) add(pattern, method string, handler Handler, filters ...Filter) {
	if !n.root {
		panic("only root can add node")
	}
	if pattern[0] != '/' {
		panic("path must begin with '/'")
	}
	patternSplitArr := strings.Split(pattern, "/")[1:]                  // [a, b, :c, d, :e, :f, g]
	n.addFunc(pattern, method, patternSplitArr, 0, handler, filters...) // 默认splitArr从0开始解析
}

// addSplitArr
//
// pattern /a/b/:c/d/:e/:f/g
//
// method eg:http.MethodGet
//
// patternSplitArr [a, b, ?, d, ?, ?, g]
//
// index 1
func (n *node) addSplitArr(pattern, method string, patternSplitArr []string, index int, handler Handler, filters ...Filter) {
	if len(patternSplitArr) == index { // splitArr长度与index相同则表明当前结点是叶子结点
		n.fill(pattern, method, handler, filters...)
		return
	}
	n.addFunc(pattern, method, patternSplitArr, index, handler, filters...)
}

// addFunc
//
// pattern /a/b/:c/d/:e/:f/g
//
// method eg:http.MethodGet
//
// patternSplitArr [a, b, ?, d, ?, ?, g]
//
// index 1
func (n *node) addFunc(pattern, method string, patternSplitArr []string, index int, handler Handler, filters ...Filter) {
	var patternPiece string
	if patternPiece = patternSplitArr[index]; patternPiece[0] == ':' {
		patternPiece = "?"
	}
	index++
	for _, nd := range n.nextNodes {
		if nd.patternPiece == patternPiece {
			if gnomon.String().IsNotEmpty(nd.method) && nd.method != method {
				break
			} else {
				nd.addSplitArr(pattern, method, patternSplitArr, index, handler, filters...)
			}
			return
		}
	}
	nextNode := nextNode(n, patternPiece)
	n.lockNode.Lock()
	n.nextNodes = append(n.nextNodes, nextNode)
	n.lockNode.Unlock()
	nextNode.addSplitArr(pattern, method, patternSplitArr, index, handler, filters...)
}

// fill
func (n *node) fill(pattern, method string, handler Handler, filters ...Filter) {
	if gnomon.String().IsNotEmpty(n.method) {
		return
	}
	n.filters = append(n.filters, filters...)
	n.pattern = pattern
	n.method = method
	n.handler = handler
	if gnomon.String().IsNotEmpty(method) {
		fmt.Printf("grope url %s %s \n", method, pattern)
	}
}

// fetch
//
// pattern /a/b/:c/d/:e/:f/g
//
// method eg:http.MethodGet
func (n *node) fetch(pattern, method string) *node {
	if !n.root {
		panic("only root can fetch node")
	}
	if pattern[0] != '/' {
		panic("path must begin with '/'")
	}
	patternSplitArr := strings.Split(pattern, "/")[1:]          // [a, b, :c, d, :e, :f, g]
	return n.fetchSplitArr(pattern, method, patternSplitArr, 0) // 默认splitArr从0开始解析
}

// fetchFunc
//
// pattern /a/b/:c/d/:e/:f/g
//
// method eg:http.MethodGet
//
// patternSplitArr [a, b, ?, d, ?, ?, g]
//
// index 1
func (n *node) fetchSplitArr(pattern string, method string, patternSplitArr []string, index int) *node {
	patternPiece := patternSplitArr[index]
	index++
	nChanStaticPiece := make(chan *node)
	nChanDynamicPiece := make(chan *node)
	count := 2
	go func() {
		nChanStaticPiece <- n.fetchFuncAsync(pattern, patternPiece, method, patternSplitArr, index)
	}()
	go func() {
		nChanDynamicPiece <- n.fetchFuncAsync(pattern, "?", method, patternSplitArr, index)
	}()
	for {
		select {
		case nd := <-nChanStaticPiece:
			count--
			if nil != nd {
				return nd
			} else if count == 0 {
				return nil
			}
		case nd := <-nChanDynamicPiece:
			count--
			if nil != nd {
				return nd
			} else if count == 0 {
				return nil
			}
		}
	}
}

func (n *node) fetchFuncAsync(pattern, patternPiece string, method string, patternSplitArr []string, index int) *node {
	for _, nd := range n.nextNodes {
		if nd.patternPiece == patternPiece {
			if len(patternSplitArr) == index { // splitArr长度与index相同则表明当前结点是叶子结点
				if nd.method == method {
					return nd
				}
				continue
			}
			return nd.fetchSplitArr(pattern, method, patternSplitArr, index)
		}
	}
	return nil
}
