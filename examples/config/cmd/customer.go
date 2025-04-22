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

	"github.com/spf13/cobra"

	config "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var customerCmd = &cobra.Command{
	Use:   "customer",
	Short: "Customer operations",
	Long:  `General commands for managing Customer`,
}

var readCustomerIDCmd = &cobra.Command{
	Use:   "readID",
	Short: "Read Customer by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var customerID string

		fmt.Print("Enter Customer ID in gid format: ")
		if _, err := fmt.Scanln(&customerID); err != nil {
			fmt.Println("Error reading customerID:", err)
		}

		resp, err := client.ReadCustomer(context.Background(), &config.ReadCustomerRequest{
			Identifier: &config.ReadCustomerRequest_Id{Id: customerID},
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readCustomerNameCmd = &cobra.Command{
	Use:   "readName",
	Short: "Read Customer by Name",
	Run: func(cmd *cobra.Command, args []string) {
		var customerName string

		fmt.Print("Enter Customer Name: ")
		if _, err := fmt.Scanln(&customerName); err != nil {
			fmt.Println("Error reading customerName:", err)
		}

		resp, err := client.ReadCustomer(context.Background(), &config.ReadCustomerRequest{
			Identifier: &config.ReadCustomerRequest_Name{Name: customerName},
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readCustomerCmd = &cobra.Command{
	Use:   "read",
	Short: "Read Customer from service account",
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := client.ReadCustomer(context.Background(), &config.ReadCustomerRequest{})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(customerCmd)
	customerCmd.AddCommand(readCustomerIDCmd)
	customerCmd.AddCommand(readCustomerNameCmd)
	customerCmd.AddCommand(readCustomerCmd)
}
