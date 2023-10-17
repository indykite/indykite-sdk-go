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

package test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"github.com/indykite/indykite-sdk-go/authorization"
	"github.com/indykite/indykite-sdk-go/grpc"
	apicfg "github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/ingest"
)

var (
	clientAuthorization *authorization.Client
	clientIngest        *ingest.Client
	retryIngest         *ingest.RetryClient
	err                 error
)

// InitConfigAuthorization reads in config file and ENV variables if set.
func InitConfigAuthorization() (*authorization.Client, error) {
	clientAuthorization, err = authorization.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)
	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Authorization Client: %v", err))
	}
	return clientAuthorization, err
}

// InitConfigIngest reads in ingest file and ENV variables if set.
func InitConfigIngest() (*ingest.Client, error) {
	clientIngest, err = ingest.NewClient(context.Background(),
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)

	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Ingest Client: %v", err))
	}
	return clientIngest, nil
}

// InitConfigIngestRetry reads in ingest file and ENV variables if set.
func InitConfigIngestRetry() (*ingest.RetryClient, error) {
	retryIngest, err = ingest.NewRetryClient(context.Background(),
		&ingest.RetryPolicy{
			MaxAttempts:       4,
			InitialBackoff:    1 * time.Second,
			BackoffMultiplier: 2,
		},
		grpc.WithCredentialsLoader(apicfg.DefaultEnvironmentLoader),
		grpc.WithRetryOptions(retry.Disable()),
	)

	if err != nil {
		er(fmt.Sprintf("failed to create IndyKite Ingest RetryClient: %v", err))
	}
	return retryIngest, nil
}

func er(msg any) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
}
