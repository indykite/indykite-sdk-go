// Copyright (c) 2022 IndyKite
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
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var invitationCmd = &cobra.Command{
	Use:   "invitation",
	Short: "Create, check, approve and check invitation",
	Run: func(cmd *cobra.Command, args []string) {

		ref := "my-random-reference-3"
		err := client.CreateEmailInvitation(context.Background(),
			"example@indykite.com",
			"696e6479-6b69-4465-8000-030f00000001",
			ref,
			time.Now().Add(time.Hour*24),
			time.Now().Add(time.Second*15),
			nil,
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client: %v", err)
		}
		fmt.Println("Invitation created")

		invState, err := client.CheckInvitationState(context.Background(), ref)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Invitation state:")
		fmt.Println(jsonp.Format(invState))

		fmt.Println("Waiting 20s")
		time.Sleep(time.Second * 20)

		invState, err = client.CheckInvitationState(context.Background(), ref)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Invitation state:")
		fmt.Println(jsonp.Format(invState))

		fmt.Println("\nWrite a token received in email")
		var token string
		fmt.Scanln(&token)

		invState, err = client.CheckInvitationToke(context.Background(), token)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Invitation state:")
		fmt.Println(jsonp.Format(invState))

		fmt.Println("Waiting 5s")
		time.Sleep(time.Second * 5)
		err = client.CancelInvitation(context.Background(), ref)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Invitation cancelled")

		invState, err = client.CheckInvitationState(context.Background(), ref)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Invitation state:")
		fmt.Println(jsonp.Format(invState))
	},
}

func init() {
	rootCmd.AddCommand(invitationCmd)
}
