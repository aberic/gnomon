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

package gnomon

import (
	"gotest.tools/assert"
	"testing"
)

func TestSM2Generate(t *testing.T) {
	pri, pub, err := SM2Generate()
	assert.NilError(t, err)
	t.Log(pri)
	t.Log(pub)
}

func TestSM2GenerateBytes(t *testing.T) {
	priBytes, pubBytes, err := SM2GenerateBytes("SM2 PRIVATE KEY", "SM2 PUBLIC KEY")
	assert.NilError(t, err)
	t.Log(string(priBytes))
	t.Log(string(pubBytes))
}
