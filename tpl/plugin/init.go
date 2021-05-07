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

package plugin

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "plugin"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Open,
			nil,
			[][2]string{
				{`{{ (plugin.Open "hello").Lookup "HelloFmt" }}`, "map[english:Hello %s french:Salutation %s spanish:Hola %s]"},
				{`{{ plugin.Call ((plugin.Open "hello").Lookup "Hello") "holyhope" }}`, "Hello holyhope"},
			},
		)

		ns.AddMethodMapping(ctx.Call,
			nil,
			[][2]string{
				{`{{ plugin.Call (plugin.Get "hello" "Hello") "holyhope" }}`, "Hello holyhope"},
				{`{{ plugin.Call ((plugin.Open "hello").Lookup "Hello") "holyhope" }}`, "Hello holyhope"},
			},
		)

		ns.AddMethodMapping(ctx.Get,
			nil,
			[][2]string{
				{`{{ plugin.Get "hello" "HelloFmt" }}`, "map[english:Hello %s french:Salutation %s spanish:Hola %s]"},
				{`{{ plugin.Call (plugin.Get "hello" "Hello") "holyhope" }}`, "Hello holyhope"},
			},
		)

		ns.AddMethodMapping(ctx.Has,
			nil,
			[][2]string{
				{`{{ plugin.Has "hello" "HelloFmt" }}`, "true"},
				{`{{ plugin.Has "hello" "Hello" }}`, "true"},
				{`{{ plugin.Has "hello" "Language" }}{{/* Constant cannot be retrieved */}}`, "false"},
				{`{{ plugin.Has "hello" "notFound" }}`, "false"},
				{`{{ plugin.Has "hello" "NotFound" }}`, "false"},
			},
		)

		ns.AddMethodMapping(ctx.Exist,
			nil,
			[][2]string{
				{`{{ "hello" | plugin.Exist }}`, "true"},
				{`{{ "does-not-exists" | plugin.Exist }}`, "false"},
				{`{{ "Hello" | plugin.Exist }}`, "false"},
			},
		)

		return ns
	}
	internal.AddTemplateFuncsNamespace(f)
}
