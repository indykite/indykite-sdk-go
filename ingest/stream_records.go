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

	ingestpb "github.com/indykite/jarvis-sdk-go/gen/indykite/ingest/v1beta1"

	"github.com/indykite/jarvis-sdk-go/errors"
)

func (c *Client) StreamRecords(mappingID string, records []*ingestpb.Record) (
	*ingestpb.StreamRecordsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	stream, err := c.client.StreamRecords(ctx)
	if err != nil {
		return nil, errors.FromError(err)
	}

	for _, record := range records {
		req := &ingestpb.StreamRecordsRequest{
			MappingConfigId: mappingID,
			Record:          record,
		}

		err = stream.Send(req)
		if err != nil {
			return nil, errors.FromError(err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.FromError(err)
	}

	log.Printf("%d records ingested succesfully", res.GetNumRecords())
	return res, nil
}
