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
	"reflect"
	"testing"
)

func TestGetPlugin(t *testing.T) {
	pluginsMu.Lock()
	plugins = map[string]PluginInfo{
		"a": {ID: "a"},
	}
	pluginsMu.Unlock()

	for i, tc := range []struct {
		input  string
		expect PluginInfo
	}{
		{
			input: "a",
			expect: PluginInfo{
				ID: "a",
			},
		},
	} {
		actual, _ := GetPlugin(tc.input)
		if !reflect.DeepEqual(actual, tc.expect) {
			t.Errorf("Test %d: Expected %v but got %v", i, tc.expect, actual)
		}
	}

}
