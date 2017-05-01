// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/v4run/bob/server"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run your project",
	Long:  `bob run allows you to manually rerun your project any time`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := http.PostForm(serverAddr(server.RUN), url.Values{"id": args})
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
