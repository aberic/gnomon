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
	"fmt"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	fmt.Println("HELLO =", Env().GetEnv("HELLO"))
}

func TestGetEnvDefault(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	fmt.Println("HELLO =", Env().GetEnvDefault("HELLO", "my"))
	fmt.Println("WORLD =", Env().GetEnvDefault("WORLD", "god"))
}

func TestGetEnvBool(t *testing.T) {
	_ = os.Setenv("HELLO", "true")
	fmt.Println("HELLO =", Env().GetEnvBool("HELLO"))
	_ = os.Setenv("HELLO", "false")
	fmt.Println("HELLO =", Env().GetEnvBool("HELLO"))
}

func TestGetEnvBoolDefault(t *testing.T) {
	fmt.Println("HELLO =", Env().GetEnvBoolDefault("HELLO", true))
	_ = os.Setenv("HELLO", "false")
	fmt.Println("HELLO =", Env().GetEnvBoolDefault("HELLO", true))
}
