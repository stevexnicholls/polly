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

package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/stevexnicholls/polly"
	log "github.com/stevexnicholls/polly/logger"
)

var (
	// Recursive _
	Recursive bool

	// Filter _
	Filter string

	// ContentDir _
	ContentDir string

	// LayoutsDir _
	LayoutsDir string

	// PollyDir _
	PollyDir string
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
	processCmd.Flags().BoolVarP(&Recursive, "recursive", "r", true, "Help message for recursive")
	processCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Help message for filter")

	processCmd.Flags().StringVarP(&ContentDir, "contentDir", "c", "content", "filesystem path to content directory")
	processCmd.Flags().StringVarP(&LayoutsDir, "layoutsDir", "l", "layouts", "filesystem path to layout directory")
	processCmd.Flags().StringVarP(&PollyDir, "pollyDir", "p", "resources"+string(os.PathSeparator)+"polly", "filesystem path to polly directory")
}

// run is the main function for this command
func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)
	defer func() {
		cancelFunction()
	}()

	config := polly.Config{
		PollyDir:   PollyDir,
		ContentDir: ContentDir,
		LayoutsDir: LayoutsDir,
	}

	plugins, err := polly.GetPlugins()
	if err != nil {
		log.Fatalw("", "error", err.Error())
	}

	// provision all plugins
	for _, p := range plugins {
		err := p.Provision(ctxWithCancel, config)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// execute all plugins
	for _, p := range plugins {
		err := p.Execute(ctxWithCancel, config)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
