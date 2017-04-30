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
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/v4run/bob/server"
)

var cfgFile string

func serverAddr(path server.Action) string {
	u := url.URL{
		Scheme: viper.GetString("scheme"),
		Host:   viper.GetString("host") + ":" + viper.GetString("port"),
		Path:   path.String(),
	}
	return u.String()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Live reload your projects",
	Long:  `bob helps you live reload all your go projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !IsRunning() {
			jww.INFO.Println("Starting server")
			if err := StartServer(); err != nil {
				jww.ERROR.Println("Error starting bob server.", err)
			}
			jww.INFO.Println("server started")
		}
		wd, _ := os.Getwd()
		dir, _ := filepath.Abs(wd)
		for !IsRunning() {
			time.Sleep(time.Millisecond * 10)
		}
		values := url.Values{
			"path": []string{
				dir,
			},
		}
		resp, err := http.PostForm(serverAddr(server.NEWJOB), values)
		if err != nil {
			jww.ERROR.Println(err)
		}
		jww.INFO.Println(resp)
	},
}

// IsRunning returns true if the server is running
func IsRunning() bool {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	if _, err := c.Get(serverAddr(server.STATUS)); err != nil {
		return false
	}
	return true
}

// StartServer bob central server
func StartServer() error {
	return exec.Command("bob", "init").Start()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bob.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	jww.SetStdoutThreshold(jww.LevelWarn)
	viper.SetDefault("port", "8989")
	viper.SetDefault("scheme", "http")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("verbose", true)
	viper.SetConfigName(".bob")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	}
	if viper.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}
