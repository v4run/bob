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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/v4run/bob/server"
)

var (
	// ErrServiceNotRunning is thrown when the bob background service is not running
	ErrServiceNotRunning = errors.New("Service not running")
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the jobs already running",
	Long:  `list lists all the projects that bob is currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(serverAddr(server.LISTJOBS))
		if err != nil {
			jww.ERROR.Println("bob background service is not running")
			return
		}
		defer resp.Body.Close()
		var js []*server.Job
		if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
			jww.ERROR.Println(err)
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(w, "ID \t Directory")
		for _, j := range js {
			fmt.Fprintf(w, "%d \t %s\n", j.ID, j.Dir)
		}
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
