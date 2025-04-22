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

package authorization

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
)

// IsAuthorized checks if DigitalTwin can perform actions on resources.
func (c *Client) IsAuthorized(
	ctx context.Context,
	digitalTwinID *authorizationpb.DigitalTwin,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinId{
				DigitalTwinId: digitalTwinID},
		},
		Resources:   resources,
		InputParams: inputParams,
		PolicyTags:  policyTags,
	}, opts...)
}

// IsAuthorizedByToken checks if DigitalTwin, identified by access token,
// can perform actions on resources.
func (c *Client) IsAuthorizedByToken(
	ctx context.Context,
	token string,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_AccessToken{
				AccessToken: token},
		},
		Resources:   resources,
		InputParams: inputParams,
		PolicyTags:  policyTags,
	}, opts...)
}

// IsAuthorizedByProperty checks if DigitalTwin, identified by property filter,
// can perform actions on resources.
func (c *Client) IsAuthorizedByProperty(
	ctx context.Context,
	property *authorizationpb.Property,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinProperty{
				DigitalTwinProperty: property},
		},
		Resources:   resources,
		InputParams: inputParams,
		PolicyTags:  policyTags,
	}, opts...)
}

// IsAuthorizedByExternalID checks if DigitalTwin, identified by external_id,
// can perform actions on resources.
func (c *Client) IsAuthorizedByExternalID(
	ctx context.Context,
	externalID *authorizationpb.ExternalID,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	inputParams map[string]*authorizationpb.InputParam,
	policyTags []string,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_ExternalId{
				ExternalId: externalID,
			},
		},
		Resources:   resources,
		InputParams: inputParams,
		PolicyTags:  policyTags,
	}, opts...)
}

func (c *Client) IsAuthorizedWithRawRequest(
	ctx context.Context,
	req *authorizationpb.IsAuthorizedRequest,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call IsAuthorized client endpoint")
	}

	if subject, ok := req.GetSubject().GetSubject().(*authorizationpb.Subject_AccessToken); ok {
		if err := verifyTokenFormat(subject.AccessToken); err != nil {
			return nil, err
		}
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.IsAuthorized(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
