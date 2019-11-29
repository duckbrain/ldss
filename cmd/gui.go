// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// +build gui

package cmd

import (
	"log"
	"net"
	"net/http"

	"github.com/duckbrain/ldss/lib/http"
	"github.com/spf13/cobra"
	"github.com/zserge/webview"
)

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the graphical program",
	Long:  `Launch the graphical program`,
	Run: func(cmd *cobra.Command, args []string) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()
		go func() {
			// Set up your http server here
			web.Handle(lang())
			log.Fatal(http.Serve(ln, nil))
		}()
		webview.Open("LDS Scriptures - By Jonathan Duck", "http://"+ln.Addr().String(), 400, 300, true)
	},
}

func init() {
	RootCmd.AddCommand(guiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// guiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// guiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
