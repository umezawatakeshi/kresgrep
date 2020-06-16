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
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

var kindRegexp *regexp.Regexp
var nameRegexp *regexp.Regexp
var namespaceRegexp *regexp.Regexp

type resource struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
}

func match(data []byte) (bool, error) {
	var res resource
	err := yaml.Unmarshal(data, &res)
	if err != nil {
		return false, err
	}
	if res.Metadata.Namespace == "" {
		res.Metadata.Namespace = "default"
	}
	return kindRegexp.MatchString(res.Kind) &&
			nameRegexp.MatchString(res.Metadata.Name) &&
			namespaceRegexp.MatchString(res.Metadata.Namespace),
		nil
}

func searchFile(filename string) {
	if filename == "-" {
		filename = "/dev/stdin"
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open failed: %v\n", err)
		return
	}
	defer file.Close()

	reader := k8syaml.NewYAMLReader(bufio.NewReader(file))
	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Read failed: %v\n", err)
			return
		}
		matched, err := match(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "YAML unmarshal failed: %v\n", err)
		} else if matched {
			str := string(data)
			fmt.Print(str)
			if !strings.HasSuffix(str, "\n") {
				fmt.Println()
			}
			fmt.Println("---")
		}
	}
}

func main() {
	kindPattern := flag.StringP("kind", "k", "", "search pattern for kind")
	namePattern := flag.StringP("name", "a", "", "search pattern for name")
	namespacePattern := flag.StringP("namespace", "n", "", "search pattern for namespace")

	flag.Parse()

	var err error
	kindRegexp, err = regexp.Compile(*kindPattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile kind pattern: %v\n", err)
		os.Exit(1)
	}
	nameRegexp, err = regexp.Compile(*namePattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile name pattern: %v\n", err)
		os.Exit(1)
	}
	namespaceRegexp, err = regexp.Compile(*namespacePattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile namespace pattern: %v\n", err)
		os.Exit(1)
	}

	var filenames []string
	args := flag.Args()
	if len(args) == 0 {
		filenames = []string{"-"}
	} else {
		filenames = args
	}
	for _, filename := range filenames {
		searchFile(filename)
	}
}
