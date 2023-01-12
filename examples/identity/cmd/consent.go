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
	"io"
	"log"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/spf13/cobra"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
)

// planCmd represents the plan command
var consentCmd = &cobra.Command{
	Use:   "consent",
	Short: "Consent operation",
	Long: `General commands for Consent

  This is a sample only.`,
}

// checkConsentChallengeCmd represents the patch command
var checkConsentChallengeCmd = &cobra.Command{
	Use:   "check_challenge",
	Short: "Check consent challenge operation",
	Long:  `Check consent challenge and return all needed info to build consent page`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter consent challenge: ")
		var consentChallenge string
		fmt.Scanln(&consentChallenge)

		resp, err := client.CheckConsentChallenge(
			context.Background(),
			&identitypb.CheckOAuth2ConsentChallengeRequest{Challenge: consentChallenge},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// createConsentVerifier represents the patch command
var createConsentVerifier = &cobra.Command{
	Use:   "create_verifier",
	Short: "Create consent verifier",
	Long:  `Create consent verifier by sending Approval or Denial result of consent page.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter consent challenge: ")
		var consentChallenge string
		fmt.Scanln(&consentChallenge)

		req := &identitypb.CreateOAuth2ConsentVerifierRequest{ConsentChallenge: consentChallenge}
		fmt.Print("Enter 1 for Approval or 2 for Denial: ")
		var result string
		fmt.Scanln(&result)

		switch strings.TrimSpace(result) {
		case "2":
			denial := &identitypb.DenialResponse{}

			fmt.Print("Enter error of denial: ")
			fmt.Scanln(&denial.Error)
			fmt.Print("Enter error description of denial: ")
			fmt.Scanln(&denial.ErrorDescription)
			fmt.Print("Enter error hint of denial: ")
			fmt.Scanln(&denial.ErrorHint)

			fmt.Println(jsonp.Format(denial))

			req.Result = &identitypb.CreateOAuth2ConsentVerifierRequest_Denial{Denial: denial}
		default:
			approval := &identitypb.ConsentApproval{}
			for {
				fmt.Print("Enter granted scope or empty line to stop: ")
				var scope string
				fmt.Scanln(&scope)
				if scope == "" {
					break
				}
				approval.GrantScopes = append(approval.GrantScopes, scope)
			}
			req.Result = &identitypb.CreateOAuth2ConsentVerifierRequest_Approval{Approval: approval}
		}

		resp, err := client.CreateConsentVerifier(context.Background(), req, retry.WithMax(2))
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var createConsent = &cobra.Command{
	Use:   "create",
	Short: "Create consent operation",
	Long:  `Create consent and return consent id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter PiiPrincipalId (DigitalTwin ID): ")
		var piiPrincipalId string
		fmt.Scanln(&piiPrincipalId)

		fmt.Print("Enter PiiProcessorId (OAuth2 App ID): ")
		var piiProcessorId string
		fmt.Scanln(&piiProcessorId)

		resp, err := client.CreateConsent(
			context.Background(),
			&identitypb.CreateConsentRequest{
				PiiPrincipalId: piiPrincipalId,
				PiiProcessorId: piiProcessorId,
				Properties:     []string{"email"},
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var revokeConsent = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke consent operation",
	Long:  `Revoke consent for user by id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter PiiPrincipalId (DigitalTwin ID): ")
		var piiPrincipalId string
		fmt.Scanln(&piiPrincipalId)

		fmt.Print("Enter consent ID: ")
		var consentId string
		fmt.Scanln(&consentId)

		_, err := client.RevokeConsent(
			context.Background(),
			&identitypb.RevokeConsentRequest{
				PiiPrincipalId: piiPrincipalId,
				ConsentIds:     []string{consentId},
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println()
	},
}

var listConsents = &cobra.Command{
	Use:   "list",
	Short: "List consents operation",
	Long:  `List all consents for a user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter PiiPrincipalId (DigitalTwin ID): ")
		var piiPrincipalId string
		fmt.Scanln(&piiPrincipalId)

		resp, err := client.ListConsents(
			context.Background(),
			&identitypb.ListConsentsRequest{
				PiiPrincipalId: piiPrincipalId,
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}

		for {
			consent, err := resp.Recv()
			if err == io.EOF {
				fmt.Println("no more consents for a user")
				break
			}
			if err != nil {
				log.Fatalf("failed to receive consent receipt %v", err)
			}
			fmt.Println("consent receipt:")
			fmt.Println(jsonp.Format(consent.ConsentReceipt))
		}
	},
}

func init() {
	rootCmd.AddCommand(consentCmd)
	consentCmd.AddCommand(checkConsentChallengeCmd)
	consentCmd.AddCommand(createConsentVerifier)
	consentCmd.AddCommand(createConsent)
	consentCmd.AddCommand(revokeConsent)
	consentCmd.AddCommand(listConsents)
}
