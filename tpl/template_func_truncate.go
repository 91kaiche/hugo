// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"errors"
	"html"
	"html/template"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cast"
)

var (
	tagRE        = regexp.MustCompile(`^<(/)?([^ ]+?)(?:(\s*/)| .*?)?>`)
	htmlSinglets = map[string]bool{
		"br": true, "col": true, "link": true,
		"base": true, "img": true, "param": true,
		"area": true, "hr": true, "input": true,
	}
)

type openTag struct {
	name string
	pos  int
}

func truncate(a interface{}, options ...interface{}) (template.HTML, error) {
	length, err := cast.ToIntE(a)
	if err != nil {
		return "", err
	}
	var textParam interface{}
	var ellipsis template.HTML

	switch len(options) {
	case 0:
		return "", errors.New("truncate requires a length and a string")
	case 1:
		textParam = options[0]
		ellipsis = " …"
	case 2:
		textParam = options[1]
		var ok bool
		if ellipsis, ok = options[0].(template.HTML); !ok {
			s, e := cast.ToStringE(options[0])
			if e != nil {
				return "", errors.New("ellipsis must be a string")
			}
			ellipsis = template.HTML(html.EscapeString(s))
		}
	default:
		return "", errors.New("too many arguments passed to truncate")
	}
	if err != nil {
		return "", errors.New("text to truncate must be a string")
	}
	text, err := cast.ToStringE(textParam)
	if err != nil {
		return "", errors.New("text must be a string")
	}

	if html, ok := textParam.(template.HTML); ok {
		return truncateHTML(length, ellipsis, html)
	}

	if utf8.RuneCountInString(text) <= length {
		return template.HTML(html.EscapeString(text)), nil
	}

	var lastWordIndex, lastNonSpace, currentLen int
	for i, r := range text {
		currentLen++
		if unicode.IsSpace(r) {
			lastWordIndex = lastNonSpace
		} else if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			lastWordIndex = i
		} else {
			lastNonSpace = i + utf8.RuneLen(r)
		}
		if currentLen > length {
			if lastWordIndex == 0 {
				return template.HTML(html.EscapeString(text[0:i])) + ellipsis, nil
			}
			return template.HTML(html.EscapeString(text[0:lastWordIndex])) + ellipsis, nil
		}
	}

	return template.HTML(html.EscapeString(text)), nil
}

func truncateHTML(length int, ellipsis, text template.HTML) (template.HTML, error) {
	if utf8.RuneCountInString(string(text)) <= length {
		return text, nil
	}

	openTags := []openTag{}

	var lastWordIndex, lastNonSpace, currentLen, endTextPos, nextTag int
	for i, r := range text {
		if i < nextTag {
			continue
		}
		slice := string(text[i:])
		m := tagRE.FindStringSubmatchIndex(slice)
		if len(m) > 0 && m[0] == 0 {
			tagname := slice[m[4]:m[5]]
			if m[2] == -1 {
				openTags = append(openTags, openTag{name: tagname, pos: i})
			} else {
				// SGML: An end tag closes, back to the matching start tag,
				// all unclosed intervening start tags with omitted end tags
				for i, tag := range openTags {
					if tag.name == tagname {
						openTags = openTags[i:]
						break
					}
				}
			}
			nextTag = i + m[1]
			continue
		}

		currentLen++
		if unicode.IsSpace(r) {
			lastWordIndex = lastNonSpace
		} else if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			lastWordIndex = i
		} else {
			lastNonSpace = i + utf8.RuneLen(r)
		}
		if currentLen > length {
			if lastWordIndex == 0 {
				endTextPos = i
			} else {
				endTextPos = lastWordIndex
			}
			out := text[0:endTextPos] + ellipsis
			for _, tag := range openTags {
				if tag.pos > endTextPos {
					break
				}
				out += ("</" + template.HTML(tag.name) + ">")
			}

			return out, nil
		}
	}

	return text, nil
}
