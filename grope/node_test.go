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
	"net/http"
	"testing"
)

func TestNodeSupport(t *testing.T) {
	root := newNode()
	root.add("/a/b/c/d", http.MethodPost, nil, nil)
	root.add("/a/b/c/d", http.MethodPut, nil, nil)
	root.add("/a/b/c/d", http.MethodPatch, nil, nil)
	root.add("/a/b/c/d/:e", http.MethodPost, nil, nil)
	root.add("/a/b/c/d/e/f", http.MethodPost, nil, nil)
	root.add("/a/:b/c/:d/e/f/g", http.MethodPost, nil, nil)
	t.Log("=============")
	printNode(root.fetch("/a/b/c/d", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d", http.MethodPut), t)
	printNode(root.fetch("/a/b/c/d", http.MethodPatch), t)
	printNode(root.fetch("/a/b/c/d/:e", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d/e/f", http.MethodPost), t)
	printNode(root.fetch("/a/:b/c/:d/e/f/g", http.MethodPost), t)
}

func TestNodeLab(t *testing.T) {
	root := newNode()
	root.add("/a/b/c/d", http.MethodPost, nil, nil)
	root.add("/a/b/c/d", http.MethodPut, nil, nil)
	root.add("/a/b/c/d", http.MethodPatch, nil, nil)
	root.add("/a/b/c/:d", http.MethodPost, nil, nil)
	root.add("/a/b/c/d/e", http.MethodPost, nil, nil)
	root.add("/a/b/c/d/:e", http.MethodPost, nil, nil)
	root.add("/a/b/c/:d/:e", http.MethodPost, nil, nil)
	root.add("/a/b/c/d/e/f", http.MethodPost, nil, nil)
	root.add("/a/b/:c/d/:e/f", http.MethodPost, nil, nil)
	root.add("/a/:b/c/d/e/:f", http.MethodPost, nil, nil)
	root.add("/a/b/c/d/e/f/g", http.MethodPost, nil, nil)
	t.Log("=============")
	printNode(root.fetch("/a/b/c/d", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d", http.MethodPut), t)
	printNode(root.fetch("/a/b/c/d", http.MethodPatch), t)
	printNode(root.fetch("/a/b/c/:d", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d/e", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d/:e", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/:d/:e", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d/e/f", http.MethodPost), t)
	printNode(root.fetch("/a/:b/c/d/e/:f", http.MethodPost), t)
	printNode(root.fetch("/a/:b/c/d/e/:f", http.MethodPost), t)
	printNode(root.fetch("/a/b/c/d/e/f/g", http.MethodPost), t)
}

func TestNodeQA(t *testing.T) {
	root := newNode()
	root.add("/v1/company/:companyID/platforms", http.MethodGet, nil, nil)
	root.add("/v1/company/:companyID/platforms", http.MethodPost, nil, nil)
	root.add("/v1/login", http.MethodPost, nil, nil)
	root.add("/v1/company/:companyID/platforms/members/invitations", http.MethodPost, nil, nil)
	root.add("/v1/company/:companyID/platforms/:platformID", http.MethodPut, nil, nil)
	t.Log("=============")
	printNode(root.fetch("/v1/company/:companyID/platforms", http.MethodGet), t)
	printNode(root.fetch("/v1/company/:companyID/platforms", http.MethodPost), t)
	printNode(root.fetch("/v1/login", http.MethodPost), t)
	printNode(root.fetch("/v1/company/:companyID/platforms/members/invitations", http.MethodPost), t)
	printNode(root.fetch("/v1/company/1/platforms/2", http.MethodPut), t)
}

func printNode(n *node, t *testing.T) {
	if nil == n {
		t.Log("none")
	} else {
		t.Log(n.method, n.patternPiece, n.pattern)
	}
}
