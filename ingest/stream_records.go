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

package ingest

import (
	"context"
	"log"
	"time"

	"github.com/indykite/indykite-sdk-go/errors"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta1"
)

func (c *Client) StreamRecords(mappingID string, records []*ingestpb.Record) (
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
			MappingConfigId: mappingID,
			Record:          record,
		}

		err = stream.Send(req)
		if err != nil {
			return nil, errors.FromError(err)
		}
		resp, err := stream.Recv()
		if err != nil {
			return nil, errors.FromError(err)
		}
		responses[i] = resp
		if resp.GetRecordError() == nil {
			log.Printf("record %d ingested succesfully", resp.GetRecordIndex())
		} else {
			log.Printf("record %v had error:", resp.GetRecordError())
		}
		err = stream.CloseSend()
		if err != nil {
			return nil, errors.FromError(err)
		}
	}
	return responses, nil
}
