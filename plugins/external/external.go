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
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/stevexnicholls/polly"
	log "github.com/stevexnicholls/polly/logger"
	util "github.com/stevexnicholls/polly/utility"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
)

func init() {
	polly.RegisterPlugin(External{})
}

// External plugin
type External struct {
}

// PollyPlugin returns a polly.PluginInfo with plugins ID
// and New function
func (External) PollyPlugin() polly.PluginInfo {
	return polly.PluginInfo{
		ID:  "external",
		New: func() polly.Plugin { return new(External) },
	}
}

// Provision _
func (External) Provision(ctx context.Context, config polly.Config) error {

	// create shortcode directory
	pluginDir := filepath.Join(config.PollyDir, "external")
	os.MkdirAll(pluginDir, os.ModePerm)

	// TODO(Steve): moves this out
	// create polly shortcode directory
	shortcodeDir := filepath.Join(config.LayoutsDir, "shortcodes", "polly")
	os.MkdirAll(shortcodeDir, os.ModePerm)

	// create external shortcode

	shortcode := "{{ with $.Page.Params.external }}{{ $file := (printf \"" + filepath.Join(config.PollyDir) + "/%s\" (sha1 .)) }}{{ readFile $file | markdownify }}{{ end }}"
	err := ioutil.WriteFile(filepath.Join(shortcodeDir, "external.html"), []byte(shortcode), 0644)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}

// Execute is the main function for the plugin
func (External) Execute(ctx context.Context, config polly.Config) error {

	// Get all markdown files from content dir
	p := []string{config.ContentDir}
	files, err := util.GetFiles(p, true, "*.md")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		log.Infof("no files found")
		return nil
	}

	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func() {
			processFile(ctx, f, config, &wg)
		}()
	}

	wg.Wait()
	log.Infow("", "plugin", "external", "status", "complete")

	return nil

}

// processFile _
func processFile(ctx context.Context, path string, config polly.Config, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Infof("external plugin processing: %s", path)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	r := bytes.NewReader(file)

	// Use Hugo's parser to parse front matter and content
	pf, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	// from Hugo code
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
		log.Infof("external plugin: external key not found: %s", path)
		return
	}
	if _, ok := pf.FrontMatter["external"].(string); !ok {
		log.Infof("external plugin: external key is not a string: %s", path)
		return
	}
	url := os.ExpandEnv(pf.FrontMatter["external"].(string))
	if !util.IsValidURL(url) {
		log.Infof("external plugin: external value is not a valid URL: %s", path)
		return
	}

	// Download file at URL
	resp, err := util.DownloadURL(ctx, url)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	var b bytes.Buffer

	// TODO(Steve): move to another plugin
	// Re-write original frontmatter (will be sorted)
	// fm := util.FormatFromString(string(pf.FrontMatterFormat))
	// err = util.WriteMap(&b, pf.FrontMatter, fm)
	// if err != nil {
	// 	log.Errorf(err.Error())
	// 	return
	// }

	// Write new content from URL
	_, err = b.Write(resp)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// Write file

	// files are named with the sha1 of the external value
	h := sha1.New()
	h.Write([]byte(pf.FrontMatter["external"].(string)))
	n := hex.EncodeToString(h.Sum(nil))
	o := filepath.Join(config.PollyDir, "external", n)

	err = ioutil.WriteFile(o, b.Bytes(), 0644)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
}
