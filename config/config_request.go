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

package config

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/wrapperspb"

	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

type (
	// NodeRequest is a request builder.
	NodeRequest struct {
		create       *configpb.CreateConfigNodeRequest
		read         *configpb.ReadConfigNodeRequest
		update       *configpb.UpdateConfigNodeRequest
		delete       *configpb.DeleteConfigNodeRequest
		listVersions *configpb.ListConfigNodeVersionsRequest
	}
)

func (x *NodeRequest) String() string {
	switch {
	case x.create != nil:
		return fmt.Sprintf("Create %s configuration", x.create.GetName())
	case x.read != nil:
		return fmt.Sprintf("Read %s configuration", x.read.GetId())
	case x.update != nil:
		return fmt.Sprintf("Update %s configuration", x.update.GetId())
	case x.delete != nil:
		return fmt.Sprintf("Delete %s configuration", x.delete.GetId())
	case x.listVersions != nil:
		return fmt.Sprintf("ListVersions of %s configuration", x.listVersions.GetId())
	default:
		return "Invalid empty request"
	}
}

func NewCreate(name string) (*NodeRequest, error) {
	if err := IsValidName(name); err != nil {
		return nil, err
	}
	return &NodeRequest{
		create: &configpb.CreateConfigNodeRequest{
			Name: name,
		},
	}, nil
}

func NewRead(id string) (*NodeRequest, error) {
	return &NodeRequest{
		read: &configpb.ReadConfigNodeRequest{
			Id: id,
		},
	}, nil
}

func NewReadWithName(name string) (*NodeRequest, error) {
	if err := IsValidName(name); err != nil {
		return nil, err
	}
	return &NodeRequest{
		// read: &configpb.ReadConfigNodeRequest{
		// 	Identifier: &configpb.ReadConfigNodeRequest_Name{
		// 		Name: name,
		// 	},
		// },
	}, nil
}

func NewUpdate(id string) (*NodeRequest, error) {
	return &NodeRequest{
		update: &configpb.UpdateConfigNodeRequest{
			Id: id,
		},
	}, nil
}

func NewDelete(id string) (*NodeRequest, error) {
	return &NodeRequest{
		delete: &configpb.DeleteConfigNodeRequest{
			Id: id,
		},
	}, nil
}

func NewListVersions(id string) (*NodeRequest, error) {
	return &NodeRequest{
		listVersions: &configpb.ListConfigNodeVersionsRequest{
			Id: id,
		},
	}, nil
}

func (x *NodeRequest) ForLocation(id string) *NodeRequest {
	if x.create != nil {
		x.create.Location = id
	}
	return x
}

// WithPreCondition sets the expected etag to check before modify or delete.
func (x *NodeRequest) WithPreCondition(etag string) *NodeRequest {
	switch {
	case x.update != nil:
		x.update.Etag = x.optionalString(etag)
	case x.delete != nil:
		x.delete.Etag = x.optionalString(etag)
	}
	return x
}

// EmptyDisplayName removes the current displayName value.
func (x *NodeRequest) EmptyDisplayName() *NodeRequest {
	return x.WithDisplayName("")
}

// WithDisplayName sets the new displayName value.
func (x *NodeRequest) WithDisplayName(v string) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.DisplayName = x.optionalString(v)
	case x.update != nil:
		x.update.DisplayName = x.optionalString(v)
	}
	return x
}

// EmptyDescription removes the current description value.
func (x *NodeRequest) EmptyDescription() *NodeRequest {
	return x.WithDescription("")
}

// WithDescription sets the new description value.
func (x *NodeRequest) WithDescription(v string) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Description = x.optionalString(v)
	case x.update != nil:
		x.update.Description = x.optionalString(v)
	}
	return x
}

func (x *NodeRequest) WithAuthorizationPolicyConfig(v *configpb.AuthorizationPolicyConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_AuthorizationPolicyConfig{AuthorizationPolicyConfig: v}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_AuthorizationPolicyConfig{AuthorizationPolicyConfig: v}
		}
	}
	return x
}

func (x *NodeRequest) WithConsentConfig(v *configpb.ConsentConfiguration) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_ConsentConfig{ConsentConfig: v}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_ConsentConfig{ConsentConfig: v}
		}
	}
	return x
}

func (x *NodeRequest) WithEventSinkConfig(v *configpb.EventSinkConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_EventSinkConfig{
				EventSinkConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_EventSinkConfig{
				EventSinkConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithTokenIntrospectConfig(v *configpb.TokenIntrospectConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_TokenIntrospectConfig{
				TokenIntrospectConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_TokenIntrospectConfig{
				TokenIntrospectConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithIngestPipelineConfig(v *configpb.IngestPipelineConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_IngestPipelineConfig{
				IngestPipelineConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_IngestPipelineConfig{
				IngestPipelineConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithEntityMatchingPipelineConfig(v *configpb.EntityMatchingPipelineConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_EntityMatchingPipelineConfig{
				EntityMatchingPipelineConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_EntityMatchingPipelineConfig{
				EntityMatchingPipelineConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithExternalDataResolverConfig(v *configpb.ExternalDataResolverConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_ExternalDataResolverConfig{
				ExternalDataResolverConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_ExternalDataResolverConfig{
				ExternalDataResolverConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithTrustScoreProfileConfig(v *configpb.TrustScoreProfileConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_TrustScoreProfileConfig{
				TrustScoreProfileConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_TrustScoreProfileConfig{
				TrustScoreProfileConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithKnowledgeQueryConfig(v *configpb.KnowledgeQueryConfig) *NodeRequest {
	switch {
	case x.create != nil:
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_KnowledgeQueryConfig{
				KnowledgeQueryConfig: v,
			}
		}
	case x.update != nil:
		x.update.Config = nil
		if v != nil {
			x.update.Config = &configpb.UpdateConfigNodeRequest_KnowledgeQueryConfig{
				KnowledgeQueryConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithCapturePipelineConfig(v *configpb.RegisterCapturePipelineConfig) *NodeRequest {
	if x.create != nil {
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_CapturePipelineConfig{
				CapturePipelineConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithCapturePipelineTopicConfig(v *configpb.RegisterCapturePipelineTopicConfig) *NodeRequest {
	if x.create != nil {
		x.create.Config = nil
		if v != nil {
			x.create.Config = &configpb.CreateConfigNodeRequest_CapturePipelineTopicConfig{
				CapturePipelineTopicConfig: v,
			}
		}
	}
	return x
}

func (x *NodeRequest) WithVersion(version int64) *NodeRequest {
	if x.read != nil {
		x.read.Version = version
	}
	return x
}

func (*NodeRequest) optionalString(v string) *wrapperspb.StringValue {
	return wrapperspb.String(strings.TrimSpace(v))
}

func (x *NodeRequest) validate() error {
	switch {
	case x.create != nil:
		if err := x.create.Validate(); err != nil {
			return err
		}
	case x.read != nil:
		if err := x.read.Validate(); err != nil {
			return err
		}
	case x.update != nil:
		if err := x.update.Validate(); err != nil {
			return err
		}
	case x.delete != nil:
		if err := x.delete.Validate(); err != nil {
			return err
		}
	case x.listVersions != nil:
		if err := x.listVersions.Validate(); err != nil {
			return err
		}
	default:
		return errors.New("empty request")
	}
	return nil
}
