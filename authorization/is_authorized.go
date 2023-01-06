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
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/indykite/jarvis-sdk-go/errors"
	authorizationpb "github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1"
	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

const (
	externalIDProperty = "extid"
)

// IsAuthorized checks if DigitalTwin can perform actions on resources.
func (c *Client) IsAuthorized(
	ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	actions []string,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &identitypb.DigitalTwinIdentifier{Filter: &identitypb.DigitalTwinIdentifier_DigitalTwin{
			DigitalTwin: digitalTwin,
		}},
		Actions:   actions,
		Resources: resources,
	}, opts...)
}

// IsAuthorizedByToken checks if DigitalTwin, identified by access token,
// can perform actions on resources.
func (c *Client) IsAuthorizedByToken(
	ctx context.Context,
	token string,
	actions []string,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &identitypb.DigitalTwinIdentifier{Filter: &identitypb.DigitalTwinIdentifier_AccessToken{
			AccessToken: token,
		}},
		Actions:   actions,
		Resources: resources,
	}, opts...)
}

// IsAuthorizedByStringExternalID checks if DigitalTwin, identified by textual ExternalID,
// can perform actions on resources.
func (c *Client) IsAuthorizedByStringExternalID(
	ctx context.Context,
	externalID string,
	actions []string,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &identitypb.DigitalTwinIdentifier{Filter: &identitypb.DigitalTwinIdentifier_PropertyFilter{
			PropertyFilter: &identitypb.PropertyFilter{
				Type:  externalIDProperty,
				Value: objects.String(externalID),
			},
		}},
		Actions:   actions,
		Resources: resources,
	}, opts...)
}

// IsAuthorizedByNumericExternalID checks if DigitalTwin, identified by numerical ExternalID,
// can perform actions on resources.
func (c *Client) IsAuthorizedByNumericExternalID(
	ctx context.Context,
	externalID int64,
	actions []string,
	resources []*authorizationpb.IsAuthorizedRequest_Resource,
	opts ...grpc.CallOption,
) (*authorizationpb.IsAuthorizedResponse, error) {
	return c.IsAuthorizedWithRawRequest(ctx, &authorizationpb.IsAuthorizedRequest{
		Subject: &identitypb.DigitalTwinIdentifier{Filter: &identitypb.DigitalTwinIdentifier_PropertyFilter{
			PropertyFilter: &identitypb.PropertyFilter{
				Type:  externalIDProperty,
				Value: objects.Int64(externalID),
			},
		}},
		Actions:   actions,
		Resources: resources,
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

	switch sub := req.Subject.Filter.(type) {
	case *identitypb.DigitalTwinIdentifier_AccessToken:
		if err := verifyTokenFormat(sub.AccessToken); err != nil {
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

func verifyTokenFormat(bearerToken string) error {
	_, err := jwt.ParseString(bearerToken, jwt.WithValidate(true), jwt.WithAcceptableSkew(time.Second))
	if err != nil {
		return errors.NewWithCause(codes.InvalidArgument, err, "invalid token format")
	}
	return nil
}
