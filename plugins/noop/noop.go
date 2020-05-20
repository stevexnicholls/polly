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

package noop

import (
	"context"

	"github.com/stevexnicholls/polly"
	log "github.com/stevexnicholls/polly/logger"
)

func init() {
	polly.RegisterPlugin(Noop{})
}

// Noop plugin
type Noop struct {
}

// PollyPlugin returns a polly.PluginInfo with plugins ID
// and New function
func (Noop) PollyPlugin() polly.PluginInfo {
	return polly.PluginInfo{
		ID:  "noop",
		New: func() polly.Plugin { return new(Noop) },
	}
}

// Provision _
func (Noop) Provision(ctx context.Context, config polly.Config) error {
	log.Debugw("", "plugin", "noop", "method", "provision")
	return nil
}

// Execute is the main function for the plugin
func (Noop) Execute(ctx context.Context, config polly.Config) error {
	log.Debugw("", "plugin", "noop", "method", "execute")
	return nil
}
