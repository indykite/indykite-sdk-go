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
	"time"

	guuid "github.com/google/uuid"
	"github.com/pborman/uuid"
	"google.golang.org/grpc"

	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"

	"github.com/indykite/jarvis-sdk-go/errors"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta1"
)

// CheckInvitationState checks the status of invitation.
func (c *Client) CheckInvitationState(ctx context.Context,
	referenceID string, opts ...grpc.CallOption) (*identitypb.CheckInvitationStateResponse, error) {
	req := &identitypb.CheckInvitationStateRequest{Identifier: &identitypb.CheckInvitationStateRequest_ReferenceId{
		ReferenceId: referenceID,
	}}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call CheckInvitationStateRequest")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CheckInvitationState(ctx, req, opts...)
	if err != nil {
		return nil, errors.FromError(err)
	}
	return resp, nil
}

// CheckInvitationToke checks the status of invitation.
func (c *Client) CheckInvitationToke(ctx context.Context,
	invitationToken string, opts ...grpc.CallOption) (*identitypb.CheckInvitationStateResponse, error) {
	req := &identitypb.CheckInvitationStateRequest{Identifier: &identitypb.CheckInvitationStateRequest_InvitationToken{
		InvitationToken: invitationToken,
	}}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call CheckInvitationStateRequest")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CheckInvitationState(ctx, req, opts...)
	if err != nil {
		return nil, errors.FromError(err)
	}
	return resp, nil
}

// ResendInvitation send invitation token to invitee.
func (c *Client) ResendInvitation(ctx context.Context, referenceID string, opts ...grpc.CallOption) error {
	req := &identitypb.ResendInvitationRequest{ReferenceId: referenceID}
	if err := req.Validate(); err != nil {
		return errors.NewInvalidArgumentErrorWithCause(err, "unable to call ResendInvitationRequest")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	_, err := c.client.ResendInvitation(ctx, req, opts...)

	return errors.FromError(err)
}

// CancelInvitation revokes a pending invitation identified by referenceID.
func (c *Client) CancelInvitation(ctx context.Context, referenceID string, opts ...grpc.CallOption) error {
	req := &identitypb.CancelInvitationRequest{ReferenceId: referenceID}
	if err := req.Validate(); err != nil {
		return errors.NewInvalidArgumentErrorWithCause(err, "unable to call CancelInvitationRequest")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	_, err := c.client.CancelInvitation(ctx, req, opts...)

	return errors.FromError(err)
}

// CreateEmailInvitation receive all properties for digital twin.
func (c *Client) CreateEmailInvitation(ctx context.Context,
	invitee string, tenantID interface{}, referenceID string,
	expireTime, inviteAtTime time.Time, messageAttributes map[string]interface{},
	opts ...grpc.CallOption) error {
	var request = &identitypb.CreateInvitationRequest{
		ReferenceId: referenceID,
		Invitee: &identitypb.CreateInvitationRequest_Email{
			Email: invitee,
		},
		InviteAtTime: optionalTime(inviteAtTime),
		ExpireTime:   optionalTime(expireTime),
	}

	switch v := tenantID.(type) {
	case string:
		var err error
		request.TenantId, err = ParseUUID(v)
		if err != nil {
			return err
		}
	case []byte:
		if uuid.UUID(v).Variant() != guuid.RFC4122 {
			return errors.NewInvalidArgumentError("invalid UUID RFC4122 variant")
		}
		request.TenantId = v
	default:
		return errors.NewInvalidArgumentError("unsupported tenantID type")
	}

	if len(messageAttributes) > 0 {
		fields, err := objects.ToMapValue(messageAttributes)
		if err != nil {
			return errors.NewInvalidArgumentErrorWithCause(err, "unable to serialise messageAttributes")
		}
		request.MessageAttributes = &objects.MapValue{Fields: fields}
	}

	_, err := c.createInvitation(ctx, request, opts...)
	return err
}

// CreateInvitationRequest receive all properties for digital twin.
func (c *Client) createInvitation(ctx context.Context,
	req *identitypb.CreateInvitationRequest,
	opts ...grpc.CallOption,
) (*identitypb.CreateInvitationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call CreateInvitation")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateInvitation(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
