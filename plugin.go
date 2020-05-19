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

package polly

import (
	"context"
	"fmt"
	"sync"
)

// Plugin _
type Plugin interface {
	PollyPlugin() PluginInfo
	Execute(context.Context, []byte) ([]byte, error)
}

// PluginInfo _
type PluginInfo struct {
	ID  string
	New func() Plugin
}

// RegisterPlugin _
func RegisterPlugin(instance Plugin) {
	plugin := instance.PollyPlugin()

	if plugin.ID == "" {
		panic("plugin ID missing")
	}

	if plugin.New == nil {
		panic("missing PluginInfo.New")
	}
	if val := plugin.New(); val == nil {
		panic("PluginInfo.New must return a non-nil plugin instance")
	}

	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	if _, ok := plugins[string(plugin.ID)]; ok {
		panic(fmt.Sprintf("plugin already registered: %s", plugin.ID))
	}
	plugins[string(plugin.ID)] = plugin
}

// GetPlugin _
func GetPlugin(id string) (PluginInfo, error) {
	if _, ok := plugins[id]; !ok {
		panic(fmt.Sprintf("plugin not registered: %s", id))
	}
	p := plugins[id]
	return p, nil
}

var (
	plugins   = make(map[string]PluginInfo)
	pluginsMu sync.RWMutex
)
