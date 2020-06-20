/*
   Copyright (C) 2020  UMEZAWA Takeshi

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"regexp"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		CaseName         string
		Yaml             string
		KindPattern      string
		NamePattern      string
		NamespacePattern string
		LabelSelector    string
		Matched          bool
		Err              bool
	}{
		{
			"yaml parse error",
			`xxx::xxx`,
			"",
			"",
			"",
			"",
			false,
			true,
		},
		{
			"match everything",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
`,
			"",
			"",
			"",
			"",
			true,
			false,
		},
		{
			"match kind",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
`,
			"a.a",
			"",
			"",
			"",
			true,
			false,
		},
		{
			"unmatch kind",
			`
kind: xxx
metadata:
  name: bbb
  namespace: ccc
`,
			"a.a",
			"",
			"",
			"",
			false,
			false,
		},
		{
			"match name",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
`,
			"",
			"b*",
			"",
			"",
			true,
			false,
		},
		{
			"unmatch name",
			`
kind: aaa
metadata:
  name: yyy
  namespace: ccc
`,
			"",
			"bb*",
			"",
			"",
			false,
			false,
		},
		{
			"match namespace",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
`,
			"",
			"",
			"c[c]c",
			"",
			true,
			false,
		},
		{
			"unmatch namespace",
			`
kind: aaa
metadata:
  name: bbb
  namespace: zzz
`,
			"",
			"",
			"c[c]c",
			"",
			false,
			false,
		},
		{
			"match default namespace",
			`
kind: aaa
metadata:
  name: bbb
`,
			"",
			"",
			"default",
			"",
			true,
			false,
		},
		{
			"match label",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
  labels:
    release: stable
`,
			"",
			"",
			"",
			"release=stable",
			true,
			false,
		},
		{
			"unmatch label",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
  labels:
    release: stable
`,
			"",
			"",
			"",
			"release=canary",
			false,
			false,
		},
		{
			"match multiple",
			`
kind: aaa
metadata:
  name: bbb
  namespace: ccc
  labels:
    release: stable
`,
			"a.a",
			"bb*",
			"c[c]c",
			"release=stable",
			true,
			false,
		},
		{
			"unmatch multiple",
			`
kind: aaa
metadata:
  name: xxx
  namespace: ccc
`,
			"a.a",
			"bb*",
			"c[c]c",
			"",
			false,
			false,
		},
	}

	for i, c := range cases {
		var err error
		nameRegexp = regexp.MustCompile(c.NamePattern)
		namespaceRegexp = regexp.MustCompile(c.NamespacePattern)
		kindRegexp = regexp.MustCompile(c.KindPattern)
		labelSelector, err = labels.Parse(c.LabelSelector)
		matched, err := match([]byte(c.Yaml))
		if matched != c.Matched || (err != nil) != c.Err {
			t.Error("case", i, "\""+c.CaseName+"\"", "failed")
		}
	}
}
