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

package ingest

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc/codes"

	"github.com/indykite/indykite-sdk-go/errors"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
	api "github.com/indykite/indykite-sdk-go/grpc"
)

type RetryPolicy struct {
	InitialBackoff    time.Duration
	MaxAttempts       int
	BackoffMultiplier int
}

type RetryClient struct {
	*Client
	clientContext       context.Context //nolint:containedctx // we need client context here
	retryPolicy         *RetryPolicy
	isUnableToReconnect bool
}

// NewRetryClient creates a new Ingest API gRPC Client with retry functionality.
func NewRetryClient(ctx context.Context, retryPolicy *RetryPolicy, opts ...api.ClientOption) (*RetryClient, error) {
	retryClientOpts := defaultRetryClientOptions()
	client, err := NewClient(ctx, append(retryClientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	c := &RetryClient{
		Client:      client,
		retryPolicy: retryPolicy,
	}
	return c, nil
}

// NewTestRetryClient creates a new test Ingest API gRPC Client with retry functionality.
func NewTestRetryClient(client ingestpb.IngestAPIClient, retryPolicy *RetryPolicy) (*RetryClient, error) {
	c, err := NewTestClient(client)
	if err != nil {
		return nil, err
	}
	rc := &RetryClient{
		Client:      c,
		retryPolicy: retryPolicy,
	}
	return rc, nil
}

// OpenStreamClient is deprecated and will be removed: use Batch functions.
func (c *RetryClient) OpenStreamClient(ctx context.Context) error {
	c.clientContext = ctx
	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return errors.FromError(err)
	}
	c.stream = stream
	return nil
}

// SendRecord sends a record on the opened stream and retries if it fails.
func (c *RetryClient) SendRecord(record *ingestpb.Record) error {
	if c.stream == nil {
		return errors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	err := c.Client.SendRecord(record)
	if err == nil {
		return nil
	}
	return c.sendRecordWithRetry(record)
}

// ReceiveResponse returns the next response available on the stream. If an error is returned (stream is closed),
// the method will wait for a server reconnect or fail.
func (c *RetryClient) ReceiveResponse() (*ingestpb.StreamRecordsResponse, error) {
	if c.stream == nil {
		return nil, errors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	resp, err := c.Client.ReceiveResponse()
	if err == nil {
		return resp, nil
	}
	return c.receiveResponseWhenReconnected()
}

// receiveResponseWhenReconnected tries to read the next response available on the stream.
// Will return error if unable to reconnect.
func (c *RetryClient) receiveResponseWhenReconnected() (*ingestpb.StreamRecordsResponse, error) {
	for {
		if c.isUnableToReconnect {
			return nil, errors.New(codes.Unavailable, "unable to reconnect to server")
		}
		time.Sleep(1 * time.Second)
		resp, err := c.Client.ReceiveResponse()
		if err == nil {
			return resp, nil
		}
	}
}

// sendRecordWithRetry creates a new stream and tries sends a record. Will retry based on retry policy.
func (c *RetryClient) sendRecordWithRetry(record *ingestpb.Record) error {
	backoffTime := c.retryPolicy.InitialBackoff
	var err error
	for i := 1; i <= c.retryPolicy.MaxAttempts; i++ {
		log.Printf("attempting to reconnect (%d/%d) in %s...", i, c.retryPolicy.MaxAttempts, backoffTime.String())
		time.Sleep(backoffTime)
		backoffTime *= time.Duration(c.retryPolicy.BackoffMultiplier)
		err = c.OpenStreamClient(c.clientContext)
		if err != nil {
			continue
		}
		time.Sleep(500 * time.Millisecond)
		err = c.Client.SendRecord(record)
		if err == nil {
			log.Printf("restablished connection to server")
			return nil
		}
	}
	log.Printf("unable to reconnect, closing client")
	c.isUnableToReconnect = true
	if closeErr := c.Close(); closeErr != nil {
		err = closeErr
	}
	return err
}

func defaultRetryClientOptions() []api.ClientOption {
	return []api.ClientOption{
		api.WithRetryOptions(retry.Disable()),
	}
}
