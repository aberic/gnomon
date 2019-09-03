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
 */

package gnomon

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	t.Log("haha empty =", String().IsEmpty("haha"))
	t.Log("'' empty =", String().IsEmpty(""))
}

func TestIsNotEmpty(t *testing.T) {
	t.Log("haha empty =", String().IsNotEmpty("haha"))
	t.Log("'' empty =", String().IsNotEmpty(""))
}

func TestConvert(t *testing.T) {
	t.Log("uu_xx_aa =", String().Convert("uu_xx_aa"))
}

func TestRandSeq(t *testing.T) {
	t.Log("13 =", String().RandSeq(13))
	t.Log("23 =", String().RandSeq(23))
	t.Log("33 =", String().RandSeq(33))
}

func TestRandSeq16(t *testing.T) {
	t.Log("RandSeq16 =", String().RandSeq16())
}

func TestTrim(t *testing.T) {
	s := "kjsdhfj ajsd\nksjhdka sjkh"
	t.Log(s, "=", String().Trim(s))
}
