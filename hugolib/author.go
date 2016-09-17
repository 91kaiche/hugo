// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

var (
	onlyNumbersRegExp = regexp.MustCompile("^[0-9]*$")
)

// Authors is a list of all authors and their metadata.
type Authors []Author

// Get returns an author from an ID
func (a Authors) Get(id string) Author {
	for _, author := range a {
		if author.ID == id {
			return author
		}
	}
	return Author{}
}

// Author contains details about the author of a page.
type Author struct {
	ID          string
	givenName   string            // givenName OR firstName
	firstName   string            // alias for GivenName
	familyName  string            // familyName OR lastName
	lastName    string            // alias for FamilyName
	displayName string            // displayName
	thumbnail   string            // thumbnail
	image       string            // image
	shortBio    string            // shortBio
	bio         string            // bio
	email       string            // email
	social      AuthorSocial      // social
	params      map[string]string // params
	Weight      int
	languages   map[string]Author
}

func (a Author) GivenName() string {
	return a.givenName
}
func (a Author) FirstName() string {
	return a.firstName
}
func (a Author) FamilyName() string {
	return a.familyName
}
func (a Author) LastName() string {
	return a.lastName
}
func (a Author) DisplayName() string {
	return a.displayName
}
func (a Author) Thumbnail() string {
	return a.thumbnail
}
func (a Author) Image() string {
	return a.image
}
func (a Author) ShortBio() string {
	return a.shortBio
}
func (a Author) Bio() string {
	return a.bio
}
func (a Author) Email() string {
	return a.email
}
func (a Author) Social() AuthorSocial {
	return a.social
}
func (a Author) Params() map[string]string {
	return a.params
}

// AuthorSocial is a place to put social usernames per author. These are the
// standard keys that themes will expect to have available, but can be
// expanded to any others on a per site basis
// - website
// - github
// - facebook
// - twitter
// - googleplus
// - pinterest
// - instagram
// - youtube
// - linkedin
// - skype
type AuthorSocial map[string]string

// URL is a convenience function that provides the correct canonical URL
// for a specific social network given a username. If an unsupported network
// is requested, only the username is returned
func (as AuthorSocial) URL(key string) string {
	switch key {
	case "github":
		return fmt.Sprintf("https://github.com/%s", as[key])
	case "facebook":
		return fmt.Sprintf("https://www.facebook.com/%s", as[key])
	case "twitter":
		return fmt.Sprintf("https://twitter.com/%s", as[key])
	case "googleplus":
		isNumeric := onlyNumbersRegExp.Match([]byte(as[key]))
		if isNumeric {
			return fmt.Sprintf("https://plus.google.com/%s", as[key])
		}
		return fmt.Sprintf("https://plus.google.com/+%s", as[key])
	case "pinterest":
		return fmt.Sprintf("https://www.pinterest.com/%s/", as[key])
	case "instagram":
		return fmt.Sprintf("https://www.instagram.com/%s/", as[key])
	case "youtube":
		return fmt.Sprintf("https://www.youtube.com/user/%s", as[key])
	case "linkedin":
		return fmt.Sprintf("https://www.linkedin.com/in/%s", as[key])
	default:
		return as[key]
	}
}

func mapToAuthors(m map[string]interface{}) Authors {
	authors := make(Authors, len(m))
	for authorID, data := range m {
		authorMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		authors = append(authors, mapToAuthor(authorID, authorMap))
	}
	sort.Stable(authors)
	return authors
}

func mapToAuthor(id string, m map[string]interface{}) Author {
	author := Author{ID: id}
	for k, data := range m {
		switch k {
		case "givenName", "firstName":
			author.givenName = cast.ToString(data)
			author.firstName = author.givenName
		case "familyName", "lastName":
			author.familyName = cast.ToString(data)
			author.lastName = author.familyName
		case "displayName":
			author.displayName = cast.ToString(data)
		case "thumbnail":
			author.thumbnail = cast.ToString(data)
		case "image":
			author.image = cast.ToString(data)
		case "shortBio":
			author.shortBio = cast.ToString(data)
		case "bio":
			author.bio = cast.ToString(data)
		case "email":
			author.email = cast.ToString(data)
		case "social":
			author.social = normalizeSocial(cast.ToStringMapString(data))
		case "params":
			author.params = cast.ToStringMapString(data)
		case "weight":
			author.Weight = cast.ToInt(data)
		case "languages":
			if author.languages == nil {
				author.languages = make(map[string]Author)
			}
			langAuthorMap := cast.ToStringMap(data)
			for lang, m := range langAuthorMap {
				authorMap, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				author.languages[lang] = mapToAuthor(id, authorMap)
			}
		}
	}

	// set a reasonable default for DisplayName
	if author.displayName == "" {
		author.displayName = author.givenName + " " + author.familyName
	}

	return author
}

// normalizeSocial makes a naive attempt to normalize social media usernames
// and strips out extraneous characters or url info
func normalizeSocial(m map[string]string) map[string]string {
	for network, username := range m {
		username = strings.TrimSpace(username)
		username = strings.TrimSuffix(username, "/")
		strs := strings.Split(username, "/")
		username = strs[len(strs)-1]
		username = strings.TrimPrefix(username, "@")
		username = strings.TrimPrefix(username, "+")
		m[network] = username
	}
	return m
}

func (a Authors) Len() int           { return len(a) }
func (a Authors) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Authors) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
