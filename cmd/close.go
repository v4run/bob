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

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/v4run/bob/server"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "closes the entire bob service.",
	Long: `close kills all the jobs that are running and stop the
	background service.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := http.Post(serverAddr(server.STOPSERVICE), "", nil)
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
}
