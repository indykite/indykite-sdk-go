/*
 * Copyright (c) 2022 IndyKite
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cmd implements the CLI commands.
package cmd

import (
	"context"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/indykite/jarvis-sdk-go/grpc"
	"github.com/indykite/jarvis-sdk-go/grpc/config"
	"github.com/indykite/jarvis-sdk-go/identity"
)

var (
	cfgFile string
	client  *identity.Client
	jsonp   = protojson.MarshalOptions{
		Multiline:       true,
		EmitUnpopulated: true,
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		if client != nil {
			_ = client.Close()
		}
	}()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "", "config file (default is $HOME/.indykite.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".indykite" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".indykite")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	client, err = identity.NewClient(context.Background(),
		grpc.WithCredentialsLoader(config.DefaultEnvironmentLoader),
		// grpc.WithTokenSource(
		// 	oauth2.StaticTokenSource(
		// 		&oauth2.Token{
		// 			AccessToken: "[TOKEN]",
		// 			TokenType:   "Bearer",
		// 		},
		// 	),
		// ),
		grpc.WithConnectionPoolSize(2),
		grpc.WithConnectionPoolSize(4),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Config Client %v", err))
	}
}

func er(msg interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
	os.Exit(1)
}
