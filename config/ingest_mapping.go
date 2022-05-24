/*
 * Copyright (c) 2022 IndyKite
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"context"

	"google.golang.org/protobuf/types/known/wrapperspb"

	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
)

func (c *Client) CreateIngestMapping(
	ctx context.Context,
	location string,
	mapping Mapping,
) (*configpb.CreateConfigNodeResponse, error) {
	request := &configpb.CreateConfigNodeRequest{
		Location:    location,
		Name:        mapping.Name,
		DisplayName: &wrapperspb.StringValue{Value: mapping.DisplayName},
		Description: &wrapperspb.StringValue{Value: mapping.Description},
		Config: &configpb.CreateConfigNodeRequest_IngestMappingConfig{
			IngestMappingConfig: &configpb.IngestMappingConfig{
				IngestType: &configpb.IngestMappingConfig_Upsert{
					Upsert: &configpb.IngestMappingConfig_UpsertData{
						Entities: mapping.Entities,
					},
				},
			},
		},
	}
	return c.CreateConfigNode(ctx, &NodeRequest{create: request})
}

func (c *Client) UpdateIngestMapping(
	ctx context.Context,
	id string,
	mapping Mapping,
) (*configpb.UpdateConfigNodeResponse, error) {
	request := &configpb.UpdateConfigNodeRequest{
		Id: id,
		Config: &configpb.UpdateConfigNodeRequest_IngestMappingConfig{
			IngestMappingConfig: &configpb.IngestMappingConfig{
				IngestType: &configpb.IngestMappingConfig_Upsert{
					Upsert: &configpb.IngestMappingConfig_UpsertData{
						Entities: mapping.Entities,
					},
				},
			},
		},
	}
	return c.UpdateConfigNode(ctx, &NodeRequest{update: request})
}

func (c *Client) GetIngestMapping(ctx context.Context, id string) (*configpb.ReadConfigNodeResponse, error) {
	request := &configpb.ReadConfigNodeRequest{
		Id: id,
	}
	return c.ReadConfigNode(ctx, &NodeRequest{read: request})
}

func (c *Client) DeleteIngestMapping(ctx context.Context, id string) (*configpb.DeleteConfigNodeResponse, error) {
	request := &configpb.DeleteConfigNodeRequest{Id: id}
	return c.DeleteConfigNode(ctx, &NodeRequest{delete: request})
}
