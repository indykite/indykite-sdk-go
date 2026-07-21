// Copyright (c) 2026 IndyKite
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
	"context"
	"net/http"

	"github.com/indykite/indykite-sdk-go/transport"
)

const pathEventSinks = "/configs/v1/event-sinks"

// KafkaSinkConfig is the Kafka event sink provider configuration.
type KafkaSinkConfig struct {
	Topic         string   `json:"topic"`
	Username      string   `json:"username,omitempty"`
	Password      string   `json:"password,omitempty"`
	DisplayName   string   `json:"display_name,omitempty"`
	LastError     string   `json:"last_error,omitempty"`
	Brokers       []string `json:"brokers"`
	DisableTLS    bool     `json:"disable_tls,omitempty"`
	TLSSkipVerify bool     `json:"tls_skip_verify,omitempty"`
}

// PubSubSinkConfig is the Google Pub/Sub event sink provider configuration.
type PubSubSinkConfig struct {
	ProjectID       string `json:"project_id"`
	TopicName       string `json:"topic_name"`
	CredentialsJSON string `json:"credentials_json"`
	DisplayName     string `json:"display_name,omitempty"`
	// LastError is reported by the platform on reads; never set it on writes.
	LastError string `json:"last_error,omitempty"`
}

// AzureEventGridSinkConfig is the Azure Event Grid event sink provider configuration.
type AzureEventGridSinkConfig struct {
	TopicEndpoint string `json:"topic_endpoint"`
	AccessKey     string `json:"access_key"`
	DisplayName   string `json:"display_name,omitempty"`
	// LastError is reported by the platform on reads; never set it on writes.
	LastError string `json:"last_error,omitempty"`
}

// AzureServiceBusSinkConfig is the Azure Service Bus event sink provider configuration.
type AzureServiceBusSinkConfig struct {
	ConnectionString string `json:"connection_string"`
	QueueOrTopicName string `json:"queue_or_topic_name"`
	DisplayName      string `json:"display_name,omitempty"`
	// LastError is reported by the platform on reads; never set it on writes.
	LastError string `json:"last_error,omitempty"`
}

// EventSinkProvider names exactly one provider backend for an event sink.
type EventSinkProvider struct {
	Kafka           *KafkaSinkConfig           `json:"kafka,omitempty"`
	PubSub          *PubSubSinkConfig          `json:"pubsub,omitempty"`
	AzureEventGrid  *AzureEventGridSinkConfig  `json:"azure_event_grid,omitempty"`
	AzureServiceBus *AzureServiceBusSinkConfig `json:"azure_service_bus,omitempty"`
}

// EventSinkFilterKeyValue is one key/value pair of a route's context filter.
type EventSinkFilterKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// EventSinkFilter matches events by type and optional context key/values.
type EventSinkFilter struct {
	EventType       string                    `json:"event_type"`
	ContextKeyValue []EventSinkFilterKeyValue `json:"context_key_value,omitempty"`
}

// EventSinkRoute routes matching events to a provider. Routes are evaluated
// sequentially; StopProcessing short-circuits later routes.
type EventSinkRoute struct {
	ProviderID     string          `json:"provider_id"`
	DisplayName    string          `json:"display_name,omitempty"`
	Filter         EventSinkFilter `json:"event_type_key_values_filter"`
	StopProcessing bool            `json:"stop_processing,omitempty"`
}

// EventSink is an event sink configuration (project scoped).
type EventSink struct {
	Providers map[string]EventSinkProvider `json:"providers"`
	Metadata
	Versioned
	Routes           []EventSinkRoute `json:"routes"`
	IncludeCDCEvents bool             `json:"include_cdc_events"`
}

// CreateEventSink is the body to create an event sink.
type CreateEventSink struct {
	ProjectID        string                       `json:"project_id"`
	Name             string                       `json:"name"`
	DisplayName      string                       `json:"display_name,omitempty"`
	Description      string                       `json:"description,omitempty"`
	Providers        map[string]EventSinkProvider `json:"providers"`
	Routes           []EventSinkRoute             `json:"routes"`
	IncludeCDCEvents bool                         `json:"include_cdc_events,omitempty"`
}

// UpdateEventSink is the body to update an event sink.
type UpdateEventSink struct {
	DisplayName      *string                      `json:"display_name,omitempty"`
	Description      *string                      `json:"description,omitempty"`
	Providers        map[string]EventSinkProvider `json:"providers"`
	Routes           []EventSinkRoute             `json:"routes"`
	IncludeCDCEvents bool                         `json:"include_cdc_events,omitempty"`
}

// EventSinkAPI is the /configs/v1/event-sinks sub-API.
type EventSinkAPI struct {
	t *transport.Client
}

// List returns the event sinks in a project.
func (a *EventSinkAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]EventSink, error) {
	return listResource[EventSink](ctx, a.t, pathEventSinks, projectListQuery(projectID, opts))
}

// Create creates an event sink.
func (a *EventSinkAPI) Create(ctx context.Context, req *CreateEventSink) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathEventSinks, req)
}

// Read fetches one event sink by gid (or by name with WithLocation).
func (a *EventSinkAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*EventSink, error) {
	return readResource[EventSink](ctx, a.t, pathEventSinks, id, readOptsQuery(opts))
}

// Update updates an event sink, optionally guarded by an ETag.
func (a *EventSinkAPI) Update(ctx context.Context, id, etag string, req *UpdateEventSink) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathEventSinks, id), req, ifMatch(etag)...)
}

// Delete deletes an event sink, optionally guarded by an ETag.
func (a *EventSinkAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathEventSinks, id, etag)
}
