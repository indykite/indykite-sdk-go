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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/spf13/cobra"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

// emailVerCmd represents the plan command
var emailVerCmd = &cobra.Command{
	Use:   "verify_email_start",
	Short: "Starts Email Verification",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID, email string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)
		fmt.Print("Enter email address: ")
		fmt.Scanln(&email)

		// anyPb, _ := anypb.New(wrapperspb.String("test"))
		// anyObject, _ := objects.Any(anyPb)

		err := client.StartEmailVerification(context.Background(),
			&identitypb.DigitalTwin{
				TenantId: tenantID,
				Id:       digitalTwinID,
			}, email,
			&objects.MapValue{
				Fields: map[string]*objects.Value{
					"name": objects.String("Grandpa"),
					// "String":        objects.String("String"),
					// "Int64":         objects.Int64(-64),
					// "UInt64":        {Value: &objects.Value_UnsignedIntegerValue{UnsignedIntegerValue: 64}},
					// "Bool":          objects.Bool(true),
					// "Float64":       objects.Float64(6.4),
					// "Time":          {Value: &objects.Value_ValueTime{ValueTime: timestamppb.Now()}},
					// "Duration":      {Value: &objects.Value_DurationValue{DurationValue: durationpb.New(time.Duration(1) * time.Hour + time.Duration(33) * time.Minute + time.Duration(10) * time.Second + time.Duration(250) * time.Millisecond)}},
					// "Identifier":    {Value: &objects.Value_IdentifierValue{IdentifierValue: objects.FromUUID(test.WillyWonkaCredentialID.UUID())}},
					// "Bytes":         {Value: &objects.Value_BytesValue{BytesValue: []byte("somefunnyjokeaboutdinosaurs")}},
					// "GeoPointValue": {Value: &objects.Value_GeoPointValue{GeoPointValue: &latlng.LatLng{Latitude: 64, Longitude: 64.09}}},
					// "Array": {Value: &objects.Value_ArrayValue{ArrayValue: &objects.ArrayValue{Values: []*objects.Value{
					//	objects.String("item1"),
					//	objects.Int64(2),
					//	objects.Null(),
					//	nil,
					// }}}},
					// "ValueNull":  objects.Null(),
					// "ObjectNull": nil,
					// "ToMapValue": {Value: &objects.Value_MapValue{MapValue: &objects.MapValue{Fields: map[string]*objects.Value{
					//	"String":     objects.String("item1"),
					//	"Int64":      objects.Int64(2),
					//	"ValueNull":  objects.Null(),
					//	"ObjectNull": nil,
					// }}}},
					// "Any1":anyObject,
				},
			},
		)
		if err != nil {
			er(err)
		}
	},
}

// emailVerCompleteCmd represents the plan command
var emailVerCompleteCmd = &cobra.Command{
	Use:   "verify_email_complete",
	Short: "Completes Email Verification",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			fmt.Print("Enter referenceId: ")
			var input string
			fmt.Scanln(&input)
			resp, err := client.VerifyDigitalTwinEmail(context.Background(), input, retry.WithMax(2))
			if err != nil {
				log.Printf("failed to invoke operation on IndyKite Client %v\n", err)
				continue
			}
			fmt.Printf("Email confirmed for: %s\n", jsonp.Format(resp))
		}
	},
}

func init() {
	dtCmd.AddCommand(emailVerCmd)
	dtCmd.AddCommand(emailVerCompleteCmd)
}
