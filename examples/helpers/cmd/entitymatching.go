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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

// createandrunCmd represents the command for creating an entitymatching.
// config node and wait until statuses are success.
var createandrunCmd = &cobra.Command{
	Use:   "createandrun",
	Short: "Create and run in entitymatching",
	Long:  `Create and run in entitymatching`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EntityMatchingPipelineConfig{
			NodeFilter: &configpb.EntityMatchingPipelineConfig_NodeFilter{
				SourceNodeTypes: []string{"Person"},
				TargetNodeTypes: []string{"Person"},
			},
			SimilarityScoreCutoff: float32(0.95),
		}
		location := "gid:AAAAApkaja7kkkkkkkkkkkk"
		name := "like-real-config-node-name23"
		similarityScoreCutoff := float32(0.95)
		response, err := client.CreateAndRunEntityMatching(
			location,
			name,
			configuration,
			similarityScoreCutoff,
		)
		if err != nil {
			fmt.Println("CreateAndRunEntityMatching: " + err.Error())
		}

		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(createandrunCmd)
}
