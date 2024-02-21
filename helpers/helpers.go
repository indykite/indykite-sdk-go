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

package helpers

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
)

// DeleteNodes is a helper function that delete all nodes of specific type either Identities
// or Resource and return array of StreamRecordsResponse.
func (c *Client) DeleteNodes(
	ctx context.Context,
	nodeType string,
) ([]*ingestpb.StreamRecordsResponse,
	error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	nodes, err := c.ClientKnowledge.ListNodes(ctx, nodeType)
	if err != nil {
		return nil, err
	}
	records := []*ingestpb.Record{}
	for _, node := range nodes {
		caser := cases.Title(language.English)
		record := &ingestpb.Record{
			Id: uuid.NewString(),
			Operation: &ingestpb.Record_Delete{
				Delete: &ingestpb.DeleteData{
					Data: &ingestpb.DeleteData_Node{
						Node: &ingestpb.NodeMatch{
							ExternalId: node.ExternalId,
							Type:       caser.String(node.Type),
						},
					},
				},
			}, // lint:file-ignore U1000 Ignore report
		}
		records = append(records, record)
	}
	responses, err := c.ClientIngest.StreamRecords(records) //nolint: contextcheck // against StreamRecords
	if err != nil {
		return nil, err
	}
	return responses, nil
}
