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

package config_test

import (
	"context"
	"fmt"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/config"
)

// Bootstrap the runtime plane: application -> agent -> credential. The signed
// credential is returned exactly once, on Create.
func ExampleAppAgentCredentialAPI_Create() {
	ctx := context.Background()
	admin, err := indykite.NewAdminFromEnv(ctx)
	if err != nil {
		return
	}

	app, err := admin.Applications().Create(ctx, &config.CreateApplication{
		ProjectID: "gid:project", Name: "checkout-service",
	})
	if err != nil {
		return
	}

	agent, err := admin.AppAgents().Create(ctx, &config.CreateAppAgent{
		ApplicationID:  app.ID,
		Name:           "checkout-agent",
		APIPermissions: []string{config.PermissionAuthorization, config.PermissionCapture},
	})
	if err != nil {
		return
	}

	cred, err := admin.AppAgentCredentials().Create(ctx, &config.CreateAppAgentCredential{
		ApplicationAgentID: agent.ID, DisplayName: "prod-key",
	})
	if err != nil {
		return
	}
	// Persist cred.ApplicationAgentConfig now — it cannot be retrieved later.
	fmt.Println(len(cred.ApplicationAgentConfig) > 0)
}

// Route audit events to Kafka with an event sink.
func ExampleEventSinkAPI_Create() {
	ctx := context.Background()
	admin, err := indykite.NewAdminFromEnv(ctx)
	if err != nil {
		return
	}

	if _, err = admin.EventSinks().Create(ctx, &config.CreateEventSink{
		ProjectID: "gid:project",
		Name:      "audit-to-kafka",
		Providers: map[string]config.EventSinkProvider{
			"kafka-main": {Kafka: &config.KafkaSinkConfig{
				Brokers: []string{"kafka.example.com:9092"},
				Topic:   "indykite-events",
			}},
		},
		Routes: []config.EventSinkRoute{{
			ProviderID: "kafka-main",
			Filter:     config.EventSinkFilter{EventType: "indykite.audit.config.create"},
		}},
	}); err != nil {
		return
	}
}
