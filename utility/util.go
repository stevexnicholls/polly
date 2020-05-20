/*
Copyright © 2020 Steve Nicholls <stevexnicholls@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Format _
type Format string

const (
	// ORG _
	ORG Format = "org"
	// JSON _
	JSON Format = "json"
	// TOML _
	TOML Format = "toml"
	// YAML _
	YAML Format = "yaml"
	// CSV _
	CSV Format = "csv"
)

const (
	yamlDelimLf = "---\n"
	tomlDelimLf = "+++\n"
	jsonDelimLf = "\n"
)

// IsValidURL tests a string to determine if it is a well-structured url or not.
func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// DownloadURL downloads a file
func DownloadURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s bytes.Buffer
	_, err = io.Copy(&s, resp.Body)
	if err != nil {
		return nil, err
	}

	return s.Bytes(), nil
}

// WriteMap pretty-prints the map based on the specified Format
func WriteMap(b *bytes.Buffer, m map[string]interface{}, f Format) error {
	switch vv := f; vv {
	case YAML:
		m, err := yaml.Marshal(m)
		if err != nil {
			return (err)
		}
		_, err = b.Write([]byte(yamlDelimLf))
		if err != nil {
			return (err)
		}
		_, err = b.Write(append(m, []byte(yamlDelimLf)...))
		if err != nil {
			return (err)
		}
	case JSON:
		m, err := json.MarshalIndent(m, "", "")
		if err != nil {
			return (err)
		}
		_, err = b.Write(append(m, []byte(jsonDelimLf)...))
		if err != nil {
			return (err)
		}
	case TOML:
		_, err := b.Write([]byte(tomlDelimLf))
		if err != nil {
			return (err)
		}
		enc := toml.NewEncoder(b)
		err = enc.Encode(m)
		if err != nil {
			return (err)
		}
		_, err = b.Write([]byte(tomlDelimLf))
		if err != nil {
			return (err)
		}
	}

	return nil
}

// FormatFromString returns a Format
func FormatFromString(formatStr string) Format {
	formatStr = strings.ToLower(formatStr)
	if strings.Contains(formatStr, ".") {
		// Assume a filename
		formatStr = strings.TrimPrefix(filepath.Ext(formatStr), ".")

	}
	switch formatStr {
	case "yaml", "yml":
		return YAML
	case "json":
		return JSON
	case "toml":
		return TOML
	case "org":
		return ORG
	case "csv":
		return CSV
	}

	return ""

}

// GetFiles _
func GetFiles(paths []string, recursive bool, filter string) ([]string, error) {

	var files []string

	// validate paths exist
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, errors.New("file does not exist")
		}
	}

	if len(paths) == 1 {
		path := paths[0]
		info, _ := os.Stat(path)

		if info.IsDir() {
			// single path
			// is a directory
			if recursive {
				// do recursive
				err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && MatchFile(filter, path) {
						files = append(files, path)
					}
					return nil
				})
				if err != nil {
					return nil, err
				}
			} else {
				// no recursive
				i, err := ioutil.ReadDir(path)
				if err != nil {
					return nil, err
				}
				for _, fi := range i {
					if !fi.IsDir() && MatchFile(filter, fi.Name()) {
						files = append(files, filepath.Join(path, fi.Name()))
					}
				}
				if err != nil {
					return nil, err
				}
			}
		} else {
			// single path
			// is file
			if MatchFile(filter, path) {
				files = append(files, path)
			}
		}
	} else {
		// multiple files
		for _, path := range paths {
			if fi, _ := os.Stat(path); !fi.IsDir() {
				if MatchFile(filter, path) {
					files = append(files, path)
				}
			}
		}
	}

	return files, nil
}

// MatchFile _
func MatchFile(pattern string, file string) bool {
	if pattern != "" {
		m, err := filepath.Match(pattern, filepath.Base(file))
		if err != nil || !m {
			return false
		}
	}
	return true
}
