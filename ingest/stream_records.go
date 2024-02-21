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
	"errors"
	"io"
	"time"

	"google.golang.org/grpc/codes"

	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
)

// StreamRecords is a helper that takes a slice of records and handles opening the stream, sending the records,
// getting the responses, and closing the stream.
func (c *Client) StreamRecords(records []*ingestpb.Record) (
	[]*ingestpb.StreamRecordsResponse,
	error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return nil, sdkerrors.FromError(err)
	}

	responses := make([]*ingestpb.StreamRecordsResponse, len(records))
	for i, record := range records {
		req := &ingestpb.StreamRecordsRequest{
			Record: record,
		}

		err = stream.Send(req)
		if err != nil {
			return nil, sdkerrors.FromError(err)
		}
		var resp *ingestpb.StreamRecordsResponse
		resp, err = stream.Recv()
		if err != nil {
			return nil, sdkerrors.FromError(err)
		}
		responses[i] = resp
	}
	err = stream.CloseSend()
	if err != nil {
		return nil, sdkerrors.FromError(err)
	}
	return responses, nil
}

// OpenStreamClient opens a stream, ready to ingest records.
func (c *Client) OpenStreamClient(ctx context.Context) error {
	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return sdkerrors.FromError(err)
	}
	c.stream = stream
	return nil
}

// SendRecord sends a record on the opened stream.
func (c *Client) SendRecord(record *ingestpb.Record) error {
	if c.stream == nil {
		return sdkerrors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	req := &ingestpb.StreamRecordsRequest{
		Record: record,
	}

	err := c.stream.Send(req)
	if err == nil || errors.Is(err, io.EOF) {
		return err
	}
	return sdkerrors.FromError(err)
}

// ReceiveResponse returns the next response available on the stream.
func (c *Client) ReceiveResponse() (*ingestpb.StreamRecordsResponse, error) {
	if c.stream == nil {
		return nil, sdkerrors.New(codes.FailedPrecondition, "a stream must be opened first")
	}

	resp, err := c.stream.Recv()
	if err == nil || errors.Is(err, io.EOF) {
		return resp, err
	}
	return nil, sdkerrors.FromError(err)
}

// CloseStream closes the gRPC stream.
func (c *Client) CloseStream() error {
	if c.stream == nil {
		return sdkerrors.New(codes.FailedPrecondition, "the stream has already been closed")
	}

	err := c.stream.CloseSend()
	if err != nil {
		return sdkerrors.FromError(err)
	}
	return nil
}
