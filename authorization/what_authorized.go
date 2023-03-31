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

	"github.com/indykite/jarvis-sdk-go/errors"
	authorizationpb "github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1"
	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
)

// WhatAuthorized returns a list of resources and allowed actions for provided resource types for
// subject, identified by DigitalTwinIdentifier, can access.
func (c *Client) WhatAuthorized(
	ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	options map[string]*authorizationpb.Option,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinIdentifier{
				DigitalTwinIdentifier: &identitypb.DigitalTwinIdentifier{
					Filter: &identitypb.DigitalTwinIdentifier_DigitalTwin{DigitalTwin: digitalTwin},
				},
			},
		},
		ResourceTypes: resourceTypes,
		Options:       options,
	}, opts...)
}

// WhatAuthorizedByToken returns a list of resources and allowed actions for provided resource types for
// subject, identified by access token, can access.
func (c *Client) WhatAuthorizedByToken(
	ctx context.Context,
	token string,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	options map[string]*authorizationpb.Option,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinIdentifier{
				DigitalTwinIdentifier: &identitypb.DigitalTwinIdentifier{
					Filter: &identitypb.DigitalTwinIdentifier_AccessToken{AccessToken: token},
				},
			},
		},
		ResourceTypes: resourceTypes,
		Options:       options,
	}, opts...)
}

// WhatAuthorizedByProperty returns a list of resources and allowed actions for provided resource types for
// subject, identified by property filter, can access.
func (c *Client) WhatAuthorizedByProperty(
	ctx context.Context,
	propertyFilter *identitypb.PropertyFilter,
	resourceTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
	options map[string]*authorizationpb.Option,
	opts ...grpc.CallOption,
) (*authorizationpb.WhatAuthorizedResponse, error) {
	return c.WhatAuthorizedWithRawRequest(ctx, &authorizationpb.WhatAuthorizedRequest{
		Subject: &authorizationpb.Subject{
			Subject: &authorizationpb.Subject_DigitalTwinIdentifier{
				DigitalTwinIdentifier: &identitypb.DigitalTwinIdentifier{
					Filter: &identitypb.DigitalTwinIdentifier_PropertyFilter{PropertyFilter: propertyFilter},
				}},
		},
		ResourceTypes: resourceTypes,
		Options:       options,
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

	if sub, ok := req.GetSubject().Subject.(*authorizationpb.Subject_DigitalTwinIdentifier); ok {
		if filter, ok := sub.DigitalTwinIdentifier.Filter.(*identitypb.DigitalTwinIdentifier_AccessToken); ok {
			if err := verifyTokenFormat(filter.AccessToken); err != nil {
				return nil, err
			}
		}
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.WhatAuthorized(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
