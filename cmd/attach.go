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
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/v4run/bob/server"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "attach to a job",
	Long:  `attach current console to the job`,
	Run: func(cmd *cobra.Command, args []string) {
		u := url.URL{
			Scheme:     "ws",
			Host:       viper.GetString("host") + ":" + viper.GetString("port"),
			Path:       server.ATTACH.String(),
			ForceQuery: true,
			RawQuery: url.Values{
				"id": args,
			}.Encode(),
		}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			jww.ERROR.Println("Couldn't connect to bob background service", err)
			return
		}
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				jww.ERROR.Println("read:", err)
				return
			}
			fmt.Printf("%s\n", message)
		}
	},
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
