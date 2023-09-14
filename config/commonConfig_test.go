// Copyright 2020 The Hugo Authors. All rights reserved.
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

package config

import (
	"errors"
	"testing"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/types"

	qt "github.com/frankban/quicktest"
)

func TestBuild(t *testing.T) {
	c := qt.New(t)

	v := New()
	v.Set("build", map[string]any{
		"useResourceCacheWhen": "always",
	})

	b := DecodeBuildConfig(v)

	c.Assert(b.UseResourceCacheWhen, qt.Equals, "always")

	v.Set("build", map[string]any{
		"useResourceCacheWhen": "foo",
	})

	b = DecodeBuildConfig(v)

	c.Assert(b.UseResourceCacheWhen, qt.Equals, "fallback")

	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.IsTrue)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.IsFalse)

	b.UseResourceCacheWhen = "always"
	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.IsTrue)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.IsTrue)
	c.Assert(b.UseResourceCache(nil), qt.IsTrue)

	b.UseResourceCacheWhen = "never"
	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.IsFalse)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.IsFalse)
	c.Assert(b.UseResourceCache(nil), qt.IsFalse)
}

func TestServer(t *testing.T) {
	c := qt.New(t)

	cfg, err := FromConfigString(`[[server.headers]]
for = "/*.jpg"

[server.headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"

[[server.redirects]]
from = "/foo/**"
to = "/foo/index.html"
status = 200

[[server.redirects]]
from = "/google/**"
to = "https://google.com/"
status = 301

[[server.redirects]]
from = "/**"
to = "/default/index.html"
status = 301



`, "toml")

	c.Assert(err, qt.IsNil)

	s, err := DecodeServer(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(s.CompileConfig(loggers.NewDefault()), qt.IsNil)

	c.Assert(s.MatchHeaders("/foo.jpg"), qt.DeepEquals, []types.KeyValueStr{
		{Key: "X-Content-Type-Options", Value: "nosniff"},
		{Key: "X-Frame-Options", Value: "DENY"},
		{Key: "X-XSS-Protection", Value: "1; mode=block"},
	})

	c.Assert(s.MatchRedirect("/foo/bar/baz"), qt.DeepEquals, Redirect{
		From:   "/foo/**",
		To:     "/foo/",
		Status: 200,
	})

	c.Assert(s.MatchRedirect("/someother"), qt.DeepEquals, Redirect{
		From:   "/**",
		To:     "/default/",
		Status: 301,
	})

	c.Assert(s.MatchRedirect("/google/foo"), qt.DeepEquals, Redirect{
		From:   "/google/**",
		To:     "https://google.com/",
		Status: 301,
	})

	// No redirect loop, please.
	c.Assert(s.MatchRedirect("/default/index.html"), qt.DeepEquals, Redirect{})
	c.Assert(s.MatchRedirect("/default/"), qt.DeepEquals, Redirect{})

	for _, errorCase := range []string{
		`[[server.redirects]]
from = "/**"
to = "/file"
status = 301`,
		`[[server.redirects]]
from = "/**"
to = "/foo/file.html"
status = 301`,
	} {

		cfg, err := FromConfigString(errorCase, "toml")
		c.Assert(err, qt.IsNil)
		_, err = DecodeServer(cfg)
		c.Assert(err, qt.IsNotNil)

	}
}

func TestBuildConfigCacheBusters(t *testing.T) {
	c := qt.New(t)
	cfg := New()
	conf := DecodeBuildConfig(cfg)
	l := loggers.NewDefault()
	c.Assert(conf.CompileConfig(l), qt.IsNil)

	m, err := conf.MatchCacheBuster(l, "assets/foo/main.js")
	c.Assert(err, qt.IsNil)
	c.Assert(m, qt.IsNotNil)
	c.Assert(m("scripts"), qt.IsTrue)
	c.Assert(m("asdf"), qt.IsFalse)

	m, _ = conf.MatchCacheBuster(l, "tailwind.config.js")
	c.Assert(m("css"), qt.IsTrue)
	c.Assert(m("js"), qt.IsFalse)

	m, err = conf.MatchCacheBuster(l, "assets/foo.json")
	c.Assert(err, qt.IsNil)
	c.Assert(m, qt.IsNotNil)
	c.Assert(m("json"), qt.IsTrue)

}

func TestBuildConfigCacheBusterstTailwindSetup(t *testing.T) {
	c := qt.New(t)
	cfg := New()
	cfg.Set("build", map[string]interface{}{
		"cacheBusters": []map[string]string{
			{
				"source": "assets/watching/hugo_stats\\.json",
				"target": "css",
			},
			{
				"source": "(postcss|tailwind)\\.config\\.js",
				"target": "css",
			},
			{
				"source": "assets/.*\\.(js|ts|jsx|tsx)",
				"target": "js",
			},
			{
				"source": "assets/.*\\.(.*)$",
				"target": "$1",
			},
		},
	})

	conf := DecodeBuildConfig(cfg)
	l := loggers.NewDefault()
	c.Assert(conf.CompileConfig(l), qt.IsNil)

	m, err := conf.MatchCacheBuster(l, "assets/watching/hugo_stats.json")
	c.Assert(err, qt.IsNil)
	c.Assert(m("css"), qt.IsTrue)
}
