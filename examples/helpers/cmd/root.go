// Copyright (c) 2023 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cmd implements the CLI commands.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/grpc"
	apicfg "github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/helpers"
	"github.com/indykite/indykite-sdk-go/ingest"
	"github.com/indykite/indykite-sdk-go/knowledge"
)

var (
	cfgFile string
	client  helpers.Client
	jsonp   = protojson.MarshalOptions{
		Multiline:       true,
		EmitUnpopulated: true,
	}
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Identity Knowledge API examples",
	Long:  `Examples of using the Identity Knowledge API to interact with data in the IndyKite Identity Knowledge Graph`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "", "config file (default is $HOME/.indykite.yaml)")
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

	client.ClientKnowledge, err = knowledge.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Identity Knowledge Client: %v", err))
	}
	client.ClientIngest, err = ingest.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Ingest Client: %v", err))
	}

	client.ClientConfig, err = config.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoaderConfig),
		grpc.WithServiceAccount(),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Config Client: %v", err))
	}

	client.ClientEntitymatching, err = entitymatching.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Entitymatching Client: %v", err))
	}
}

func er(msg any) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
	os.Exit(1)
}
