// Copyright 2019 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"github.com/gohugoio/hugo/identity"
)

// Increments on breaking changes.
const TemplateVersion = 2

// Info holds some info extracted from a parsed template.
type Info struct {

	// Set for shortcode templates with any {{ .Inner }}
	IsInner bool

	// Set for partials with a return statement.
	HasReturn bool

	// Config extracted from template.
	Config Config

	identity.Manager
}

func (info Info) IsZero() bool {
	return info.Config.Version == 0
}

type Config struct {
	Version int
}

var DefaultConfig = Config{
	Version: TemplateVersion,
}
