/*
 *  Copyright (c) 2020. aberic - All Rights Reserved.
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

package log

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// debugLevel logs are typically voluminous, and are usually disabled in
	// production.
	debugLevel Level = iota - 1
	// infoLevel is the default logging priority.
	infoLevel
	// warnLevel logs are more important than Info, but don't need individual
	// human review.
	warnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	errorLevel
	// panicLevel logs a message, then panics.
	panicLevel
	// fatalLevel logs a message, then calls os.Exit(1).
	fatalLevel
	allLevel
)

const (
	logNameDebug = "DEBUG"
	logNameInfo  = "INFO "
	logNameWarn  = "WARN "
	logNameError = "ERROR"
	logNamePanic = "PANIC"
	logNameFatal = "FATAL"
)
