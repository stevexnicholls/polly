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
	"io/ioutil"
	"sync"

	"github.com/spf13/cobra"

	"github.com/stevexnicholls/polly"
	log "github.com/stevexnicholls/polly/logger"
	util "github.com/stevexnicholls/polly/utility"
)

// Recursive _
var Recursive bool

// Filter _
var Filter string

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
	processCmd.Flags().BoolVarP(&Recursive, "recursive", "r", false, "Help message for recursive")
	processCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Help message for filter")
}

// run is the main function for this command
func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)
	defer func() {
		cancelFunction()
	}()

	files, err := util.GetFiles(args, Recursive, Filter)
	if err != nil {
		log.Errorf(err.Error())
	}

	// TODO: this is just for dev; replace all of this
	ext, err := polly.GetPlugin("external")
	if err != nil {
		log.Fatalf(err.Error())
	}
	exti := ext.New()
	if err != nil {
		log.Fatalf(err.Error())
	}
	var stack []polly.Plugin
	stack = append(stack, exti)

	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go processFile(ctxWithCancel, f, stack, &wg)
	}

	wg.Wait()

	log.Infof("done")
}

// processFile process a file using available plugins
func processFile(ctx context.Context, path string, stack []polly.Plugin, wg *sync.WaitGroup) error {
	defer wg.Done()

	log.Infof("process %s", path)

	var err error

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, p := range stack {
		content, err = p.Execute(ctx, content)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}
