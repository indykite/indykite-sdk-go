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

package identity

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/indykite/jarvis-sdk-go/errors"
	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta1"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

// StartEmailVerification function initiates the flow where Indykite systems sends a
// notification to DigitalTwin with a link to verify the control over
// the notification channel (email only for now).
//
// This is a protected operation and it can be accessed only with valid agent credentials!
func (c *Client) StartEmailVerification(ctx context.Context, digitalTwin *identitypb.DigitalTwin, email string,
	attributes *objects.MapValue, opts ...grpc.CallOption) error {
	if _, _, err := digitalTwin.Verify(); err != nil {
		return errors.NewInvalidArgumentError(err.Error())
	}
	req := &identitypb.StartDigitalTwinEmailVerificationRequest{
		DigitalTwin: digitalTwin,
		Email:       email,
		Attributes:  attributes,
	}
	if err := req.Validate(); err != nil {
		return errors.NewInvalidArgumentErrorWithCause(err, "unable to call StartEmailVerification")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	_, err := c.client.StartDigitalTwinEmailVerification(ctx, req, opts...)

	return errors.FromError(err)
}

// VerifyDigitalTwinEmail function confirms to IndyKite system that the message from
// StartDigitalTwinEmailVerification function was sent and user visited the link.
//
// This is a protected operation and it can be accessed only with valid agent credentials!
func (c *Client) VerifyDigitalTwinEmail(ctx context.Context,
	token string, opts ...grpc.CallOption) (*identitypb.DigitalTwin, error) {
	req := &identitypb.VerifyDigitalTwinEmailRequest{
		Token: token,
	}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call VerifyDigitalTwinEmail")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.VerifyDigitalTwinEmail(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp.GetDigitalTwin(), nil
}

// IntrospectToken function validates the token and returns information about it.
//
// This is a protected operation and it can be accessed only with valid agent credentials!
func (c *Client) IntrospectToken(ctx context.Context,
	token string, opts ...grpc.CallOption) (*identitypb.TokenIntrospectResponse, error) {
	if err := verifyTokenFormat(token); err != nil {
		return &identitypb.TokenIntrospectResponse{Active: false}, nil
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.TokenIntrospect(ctx, &identitypb.TokenIntrospectRequest{Token: token}, opts...)

	switch s := errors.FromError(err); {
	case s == nil:
		return resp, nil
	case s.Code() == codes.InvalidArgument:
		return &identitypb.TokenIntrospectResponse{Active: false}, nil
	default:
		return nil, s
	}
}

// ChangeMyPassword change password of DigitalTwin from bearer token.
func (c *Client) ChangeMyPassword(ctx context.Context,
	token string,
	newPassword string,
	opts ...grpc.CallOption,
) error {
	if err := verifyTokenFormat(token); err != nil {
		return err
	}
	if len(newPassword) <= 4 {
		return errors.NewInvalidArgumentError("password is too short")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ChangePassword(ctx, &identitypb.ChangePasswordRequest{
		Uid:      &identitypb.ChangePasswordRequest_Token{Token: token},
		Password: newPassword,
	}, opts...)

	if s := errors.FromError(err); s != nil {
		return s
	}
	if resp.Error != nil {
		return errors.NewInvalidArgumentError(resp.Error.Code)
	}
	return nil
}

// ChangePasswordOfDigitalTwin change password of passed DigitalTwin
func (c *Client) ChangePasswordOfDigitalTwin(ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	newPassword string,
	opts ...grpc.CallOption,
) (*identitypb.ChangePasswordResponse, error) {
	if _, _, err := digitalTwin.Verify(); err != nil {
		return nil, errors.NewInvalidArgumentError(err.Error())
	}
	if len(newPassword) <= 4 {
		return nil, errors.NewInvalidArgumentError("password is too short")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ChangePassword(ctx, &identitypb.ChangePasswordRequest{
		Uid:      &identitypb.ChangePasswordRequest_DigitalTwin{DigitalTwin: digitalTwin},
		Password: newPassword,
	}, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) getDigitalTwinWithProperties(ctx context.Context,
	identifier *identitypb.DigitalTwinIdentifier,
	properties []*identitypb.PropertyMask,
	opts ...grpc.CallOption,
) (*identitypb.GetDigitalTwinResponse, error) {
	req := &identitypb.GetDigitalTwinRequest{Id: identifier, Properties: properties}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call GetDigitalTwin")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.GetDigitalTwin(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

// GetDigitalTwin receive all properties for given digital twin.
func (c *Client) GetDigitalTwin(ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	properties []*identitypb.PropertyMask,
	opts ...grpc.CallOption,
) (*identitypb.GetDigitalTwinResponse, error) {
	if _, _, err := digitalTwin.Verify(); err != nil {
		return nil, errors.NewInvalidArgumentError(err.Error())
	}

	return c.getDigitalTwinWithProperties(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_DigitalTwin{DigitalTwin: digitalTwin},
	}, properties, opts...)
}

// GetDigitalTwinByToken receive all properties for digital twin.
func (c *Client) GetDigitalTwinByToken(ctx context.Context,
	token string,
	properties []*identitypb.PropertyMask,
	opts ...grpc.CallOption,
) (*identitypb.GetDigitalTwinResponse, error) {
	if err := verifyTokenFormat(token); err != nil {
		return nil, err
	}

	return c.getDigitalTwinWithProperties(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_AccessToken{AccessToken: token},
	}, properties, opts...)
}

func (c *Client) patchDigitalTwinProperties(ctx context.Context,
	identifier *identitypb.DigitalTwinIdentifier,
	operations []*identitypb.PropertyBatchOperation,
	forceDelete bool,
	opts ...grpc.CallOption,
) ([]*identitypb.BatchOperationResult, error) {
	if err := identitypb.PropertyBatchOperations(operations).Validate(false); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "invalid patch request")
	}
	req := &identitypb.PatchDigitalTwinRequest{
		Id:          identifier,
		Operations:  operations,
		ForceDelete: forceDelete,
	}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call PatchDigitalTwin")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.PatchDigitalTwin(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp.GetResult(), nil
}

// PatchDigitalTwin update properties for given digital twin.
func (c *Client) PatchDigitalTwin(ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	operations []*identitypb.PropertyBatchOperation,
	forceDelete bool,
	opts ...grpc.CallOption,
) ([]*identitypb.BatchOperationResult, error) {
	if _, _, err := digitalTwin.Verify(); err != nil {
		return nil, errors.NewInvalidArgumentError(err.Error())
	}

	return c.patchDigitalTwinProperties(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_DigitalTwin{DigitalTwin: digitalTwin},
	}, operations, forceDelete, opts...)
}

// PatchDigitalTwinByToken update properties for digital twin.
func (c *Client) PatchDigitalTwinByToken(ctx context.Context,
	token string,
	operations []*identitypb.PropertyBatchOperation,
	forceDelete bool,
	opts ...grpc.CallOption,
) ([]*identitypb.BatchOperationResult, error) {
	if err := verifyTokenFormat(token); err != nil {
		return nil, err
	}

	return c.patchDigitalTwinProperties(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_AccessToken{AccessToken: token},
	}, operations, forceDelete, opts...)
}

func (c *Client) deleteDigitalTwin(ctx context.Context,
	identifier *identitypb.DigitalTwinIdentifier,
	opts ...grpc.CallOption,
) (*identitypb.DigitalTwin, error) {
	req := &identitypb.DeleteDigitalTwinRequest{Id: identifier}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call DeleteDigitalTwin")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteDigitalTwin(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp.GetDigitalTwin(), nil
}

// DeleteDigitalTwin deletes the digital twin.
func (c *Client) DeleteDigitalTwin(ctx context.Context,
	digitalTwin *identitypb.DigitalTwin,
	opts ...grpc.CallOption,
) (*identitypb.DigitalTwin, error) {
	if _, _, err := digitalTwin.Verify(); err != nil {
		return nil, errors.NewInvalidArgumentError(err.Error())
	}

	return c.deleteDigitalTwin(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_DigitalTwin{DigitalTwin: digitalTwin},
	}, opts...)
}

// DeleteDigitalTwinByToken deletes the digital twin.
func (c *Client) DeleteDigitalTwinByToken(ctx context.Context,
	token string,
	opts ...grpc.CallOption,
) (*identitypb.DigitalTwin, error) {
	if err := verifyTokenFormat(token); err != nil {
		return nil, err
	}

	return c.deleteDigitalTwin(ctx, &identitypb.DigitalTwinIdentifier{
		Filter: &identitypb.DigitalTwinIdentifier_AccessToken{AccessToken: token},
	}, opts...)
}
