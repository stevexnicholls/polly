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

package external

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/stevexnicholls/polly"
	log "github.com/stevexnicholls/polly/logger"
	util "github.com/stevexnicholls/polly/utility"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
)

// External plugin
type External struct {
}

func init() {
	polly.RegisterPlugin(External{})
}

// PollyPlugin returns a polly.PluginInfo with plugins ID
// and New function
func (External) PollyPlugin() polly.PluginInfo {
	return polly.PluginInfo{
		ID:  "external",
		New: func() polly.Plugin { return new(External) },
	}
}

// Execute is the main function for the plugin
func (External) Execute(ctx context.Context, b []byte) ([]byte, error) {
	var err error

	r := bytes.NewReader(b)
	pf, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// Taken from Hugo code
	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	}

	// Validate external key exists and is a valid URL
	if _, ok := pf.FrontMatter["external"]; !ok {
		log.Infof("external key not found in front matter")
		return b, nil
	}
	if _, ok := pf.FrontMatter["external"].(string); !ok {
		log.Infof("external key is not a string")
		return b, nil
	}
	u := os.ExpandEnv(pf.FrontMatter["external"].(string))
	if !util.IsValidURL(u) {
		log.Infof("external value is not a valid URL")
		return b, nil
	}

	// Download file at URL
	resp, err := util.DownloadURL(ctx, u)
	if err != nil {
		return b, err
	}

	var ret bytes.Buffer

	// Re-write original frontmatter (will be sorted)
	fmt := util.FormatFromString(string(pf.FrontMatterFormat))
	err = util.WriteMap(&ret, pf.FrontMatter, fmt)
	if err != nil {
		return b, err
	}

	// Write new content from URL
	_, err = ret.Write(resp)
	if err != nil {
		return b, err
	}

	return ret.Bytes(), nil
}
