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
	"time"

	"google.golang.org/grpc/codes"

	"github.com/indykite/jarvis-sdk-go/errors"
	ingestpb "github.com/indykite/jarvis-sdk-go/gen/indykite/ingest/v1beta2"
)

// StreamRecords is a helper that takes a slice of records and handles opening the stream, sending the records,
// getting the responses, and closing the stream.
func (c *V2Client) StreamRecords(records []*ingestpb.Record) (
	[]*ingestpb.StreamRecordsResponse,
	error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return nil, errors.FromError(err)
	}

	responses := make([]*ingestpb.StreamRecordsResponse, len(records))
	for i, record := range records {
		req := &ingestpb.StreamRecordsRequest{
			Record: record,
		}

		err = stream.Send(req)
		if err != nil {
			return nil, errors.FromError(err)
		}
		var resp *ingestpb.StreamRecordsResponse
		resp, err = stream.Recv()
		if err != nil {
			return nil, errors.FromError(err)
		}
		responses[i] = resp
	}
	err = stream.CloseSend()
	if err != nil {
		return nil, errors.FromError(err)
	}
	return responses, nil
}

// OpenStreamClient opens a stream, ready to ingest records.
func (c *V2Client) OpenStreamClient(ctx context.Context) error {
	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return errors.FromError(err)
	}
	c.stream = stream
	return nil
}

// SendRecord sends a record on the opened stream.
func (c *V2Client) SendRecord(record *ingestpb.Record) error {
	if c.stream == nil {
		return errors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	req := &ingestpb.StreamRecordsRequest{
		Record: record,
	}

	err := c.stream.Send(req)
	if err != nil {
		return errors.FromError(err)
	}
	return nil
}

// ReceiveResponse returns the next response available on the stream.
func (c *V2Client) ReceiveResponse() (*ingestpb.StreamRecordsResponse, error) {
	if c.stream == nil {
		return nil, errors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	resp, err := c.stream.Recv()
	if err != nil {
		return nil, errors.FromError(err)
	}
	return resp, nil
}

// CloseStream closes the gRPC stream.
func (c *V2Client) CloseStream() error {
	if c.stream == nil {
		return errors.New(codes.FailedPrecondition, "the stream has already been closed")
	}

	err := c.stream.CloseSend()
	if err != nil {
		return errors.FromError(err)
	}
	return nil
}
