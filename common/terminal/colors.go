// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package terminal contains helper for the terminal, such as coloring output.
package terminal

import (
	"strings"
)

const (
	errorColor   = "\033[1;31m%s\033[0m"
	warningColor = "\033[0;33m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
)

// Notice colorizes the string in a noticeable color.
func Notice(s string) string {
	return colorize(s, noticeColor)
}

// Error colorizes the string in a colour that grabs attention.
func Error(s string) string {
	return colorize(s, errorColor)
}

// Warning colorizes the string in a colour that warns.
func Warning(s string) string {
	return colorize(s, warningColor)
}

func doublePercent(str string) string {
	return strings.Replace(str, "%", "%%", -1)
}

func singlePercent(str string) string {
	return strings.Replace(str, "%%", "%", -1)
}
