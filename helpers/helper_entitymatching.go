// Copyright (c) 2024 IndyKite
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
	"time"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	entitymatchingpb "github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1"
)

type CreateAndRun struct {
	ConfigNodeID string
}

func (c *Client) CreateAndRunEntityMatching(
	location string,
	name string,
	configuration *configpb.EntityMatchingPipelineConfig,
	similarityScoreCutoff float32,
) (*CreateAndRun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	duration := time.Second * 10

	// create entitymatching config node
	createReq, _ := config.NewCreate(name)
	createReq.ForLocation(location)
	createReq.WithEntityMatchingPipelineConfig(configuration)
	resp, err := c.ClientConfig.CreateConfigNode(ctx, createReq)
	if err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "failed to Create ConfigNode on IndyKite Client")
	}
	idPipeline := resp.Id

	// call ReadSuggestedPropertyMapping with interval until success
	if !readSuggestedPropertyMappingWithInterval(ctx, c, idPipeline, duration) {
		return nil, errors.NewInvalidArgumentErrorWithCause(
			err,
			"failed to run readSuggestedPropertyMappingWithInterval on IndyKite Client")
	}

	// run RunEntityMatchingPipeline
	reqRun := &entitymatchingpb.RunEntityMatchingPipelineRequest{
		Id:                    idPipeline,
		SimilarityScoreCutoff: similarityScoreCutoff,
	}
	respRun, err := c.ClientEntitymatching.RunEntityMatchingPipeline(context.Background(), reqRun)
	if err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(
			err,
			"failed to run RunEntityMatchingPipeline on IndyKite Client")
	}

	if respRun.Id != idPipeline {
		return nil, errors.NewInvalidArgumentError(
			"failed to retrieve id from  RunEntityMatchingPipeline on IndyKite Client")
	}

	// call ReadConfigNode with interval until success
	if !readConfigNodeWithInterval(ctx, c, idPipeline, duration) {
		return nil, errors.NewInvalidArgumentErrorWithCause(
			err,
			"failed to run readConfigNodeWithInterval on IndyKite Client")
	}
	return &CreateAndRun{ConfigNodeID: idPipeline}, nil
}

func readSuggestedPropertyMappingWithInterval(ctx context.Context, c *Client, id string, interval time.Duration) bool {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			reqPropertyMapping := &entitymatchingpb.ReadSuggestedPropertyMappingRequest{
				Id: id,
			}
			//nolint:contextcheck //the context is different
			respPropertyMapping, err := c.ClientEntitymatching.ReadSuggestedPropertyMapping(
				context.Background(),
				reqPropertyMapping)
			if err != nil {
				return false
			}
			if respPropertyMapping.PropertyMappingStatus ==
				entitymatchingpb.PipelineStatus_PIPELINE_STATUS_STATUS_SUCCESS {
				return true
			}
		}
	}
}

func readConfigNodeWithInterval(ctx context.Context, c *Client, id string, interval time.Duration) bool {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			configNodeRequest, err := config.NewRead(id)
			if err != nil {
				return false
			}
			//nolint:contextcheck //the context is different
			configNodeResponse, err := c.ClientConfig.ReadConfigNode(
				context.Background(),
				configNodeRequest)
			if err != nil {
				return false
			}
			if configNodeResponse.ConfigNode.GetEntityMatchingPipelineConfig().EntityMatchingStatus ==
				configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS {
				return true
			}
		}
	}
}
