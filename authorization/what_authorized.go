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

package authorization

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
)

// WhatAuthorized returns a list of resources and allowed actions for provided resource types for
// subject, identified by DigitalTwinIdentifier, can access.
func (c *Client) WhatAuthorized(
	ctx context.Context,
	digitalTwinID *authorizationpb.DigitalTwin,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinId{
				DigitalTwinId: digitalTwinID},
		},
		ResourceTypes: resourceTypes,
		InputParams:   inputParams,
		PolicyTags:    policyTags,
	}, opts...)
}

// WhatAuthorizedByToken returns a list of resources and allowed actions for provided resource types for
// subject, identified by access token, can access.
func (c *Client) WhatAuthorizedByToken(
	ctx context.Context,
	token string,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_AccessToken{
				AccessToken: token},
		},
		ResourceTypes: resourceTypes,
		InputParams:   inputParams,
		PolicyTags:    policyTags,
	}, opts...)
}

// WhatAuthorizedByProperty returns a list of resources and allowed actions for provided resource types for
// subject, identified by property filter, can access.
func (c *Client) WhatAuthorizedByProperty(
	ctx context.Context,
	property *authorizationpb.Property,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinProperty{
				DigitalTwinProperty: property},
		},
		ResourceTypes: resourceTypes,
		InputParams:   inputParams,
		PolicyTags:    policyTags,
	}, opts...)
}

// WhatAuthorizedByExternalID returns a list of resources and allowed actions for provided resource types for
// subject, identified by external_id, can access.
func (c *Client) WhatAuthorizedByExternalID(
	ctx context.Context,
	externalID *authorizationpb.ExternalID,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_ExternalId{
				ExternalId: externalID,
			},
		},
		ResourceTypes: resourceTypes,
		InputParams:   inputParams,
		PolicyTags:    policyTags,
	}, opts...)
}

// WhatAuthorizedWithRawRequest returns a list of resources and allowed actions for provided resource types for
// subject can access.
func (c *Client) WhatAuthorizedWithRawRequest(
	ctx context.Context,
	req *authorizationpb.WhatAuthorizedRequest,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call WhatAuthorized client endpoint")
	}

	if subject, ok := req.GetSubject().GetSubject().(*authorizationpb.Subject_AccessToken); ok {
		if err := verifyTokenFormat(subject.AccessToken); err != nil {
			return nil, err
		}
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.WhatAuthorized(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
