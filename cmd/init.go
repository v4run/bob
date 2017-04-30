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
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/v4run/bob/server"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialises bob server",
	Long: `Starts the central controlling server for bob.
	It is called automatically when bob runs. User need use this manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !IsRunning() {
			s := server.New()
			if err := s.Serve(); err != nil {
				jww.ERROR.Println("Error starting bob server.", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
