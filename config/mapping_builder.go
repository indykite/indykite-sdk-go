/*
 * Copyright (c) 2022 IndyKite
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"errors"

	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
)

// Mapping Holds the values of an Ingest Mapping Config.
type Mapping struct {
	Name        string
	DisplayName string
	Description string
	Entities    []*configpb.IngestMappingConfig_Entity
}

// MappingBuilder supplies convenient functions to build an Ingest Mapping Config.
type MappingBuilder struct {
	Mapping *Mapping
}

// DigitalTwinBuilder supplies convenient functions to build a DigitalTwin entity.
type DigitalTwinBuilder struct {
	DigitalTwin *configpb.IngestMappingConfig_Entity
}

// EntityBuilder supplies convenient functions to build a domain entity.
type EntityBuilder struct {
	Entity *configpb.IngestMappingConfig_Entity
}

// NewDigitalTwinBuilder Initiates a new DigitalTwinBuilder.
func NewDigitalTwinBuilder() *DigitalTwinBuilder {
	return &DigitalTwinBuilder{DigitalTwin: &configpb.IngestMappingConfig_Entity{}}
}

// NewEntityBuilder Initiates a new EntityBuilder.
func NewEntityBuilder() *EntityBuilder {
	return &EntityBuilder{Entity: &configpb.IngestMappingConfig_Entity{}}
}

// NewMappingBuilder Initiates a new MappingBuilder.
func NewMappingBuilder() *MappingBuilder {
	return &MappingBuilder{Mapping: &Mapping{}}
}

// Name Sets the name of the Mapping.
func (m *MappingBuilder) Name(name string) *MappingBuilder {
	m.Mapping.Name = name
	return m
}

// DisplayName sets the display name of the Mapping.
func (m *MappingBuilder) DisplayName(displayName string) *MappingBuilder {
	m.Mapping.DisplayName = displayName
	return m
}

// Description sets the description of the Mapping.
func (m *MappingBuilder) Description(description string) *MappingBuilder {
	m.Mapping.Description = description
	return m
}

// DigitalTwins sets the Digital Twin entities in the Mapping, setting the label to 'DigitalTwin'.
func (m *MappingBuilder) DigitalTwins(digitalTwins []*configpb.IngestMappingConfig_Entity) *MappingBuilder {
	for _, dtb := range digitalTwins {
		dtb.Labels = []string{"DigitalTwin"}
		m.Mapping.Entities = append(m.Mapping.Entities, dtb)
	}
	return m
}

// Entities sets the non-digital twin entities in the Mapping.
func (m *MappingBuilder) Entities(entities []*configpb.IngestMappingConfig_Entity) *MappingBuilder {
	m.Mapping.Entities = append(m.Mapping.Entities, entities...)
	return m
}

// ExternalID sets the externalID property in the DigitalTwin.
func (dt *DigitalTwinBuilder) ExternalID(externalID string) *DigitalTwinBuilder {
	dt.DigitalTwin.ExternalId = &configpb.IngestMappingConfig_Property{
		SourceName: externalID,
		MappedName: "ExternalId",
		IsRequired: true,
	}
	return dt
}

// TenantID sets the tenantID in the DigitalTwin.
func (dt *DigitalTwinBuilder) TenantID(tenantID string) *DigitalTwinBuilder {
	dt.DigitalTwin.TenantId = tenantID
	return dt
}

// Properties sets the properties in the DigitalTwin.
func (dt *DigitalTwinBuilder) Properties(properties []*configpb.IngestMappingConfig_Property) *DigitalTwinBuilder {
	dt.DigitalTwin.Properties = append(dt.DigitalTwin.Properties, properties...)
	return dt
}

// Relationships sets the relationships in the DigitalTwin.
func (dt *DigitalTwinBuilder) Relationships(
	relationships []*configpb.IngestMappingConfig_Relationship) *DigitalTwinBuilder {
	dt.DigitalTwin.Relationships = append(dt.DigitalTwin.Relationships, relationships...)
	return dt
}

// Build returns the complete DigitalTwin Entity.
func (dt *DigitalTwinBuilder) Build() *configpb.IngestMappingConfig_Entity {
	return dt.DigitalTwin
}

// Build returns the complete non-digital twin Entity.
func (eb *EntityBuilder) Build() *configpb.IngestMappingConfig_Entity {
	return eb.Entity
}

// ExternalID sets the externalID property in the Entity.
func (eb *EntityBuilder) ExternalID(externalID string) *EntityBuilder {
	eb.Entity.ExternalId = &configpb.IngestMappingConfig_Property{
		SourceName: externalID,
		MappedName: "ExternalId",
		IsRequired: true,
	}
	return eb
}

// Properties sets the properties in the Entity.
func (eb *EntityBuilder) Properties(properties []*configpb.IngestMappingConfig_Property) *EntityBuilder {
	eb.Entity.Properties = append(eb.Entity.Properties, properties...)
	return eb
}

// Relationships sets the relationships in the Entity.
func (eb *EntityBuilder) Relationships(
	relationships []*configpb.IngestMappingConfig_Relationship) *EntityBuilder {
	eb.Entity.Relationships = append(eb.Entity.Relationships, relationships...)
	return eb
}

// Labels sets the relationships in the Entity.
func (eb *EntityBuilder) Labels(
	labels []string) *EntityBuilder {
	eb.Entity.Labels = append(eb.Entity.Labels, labels...)
	return eb
}

func (m *Mapping) Validate() error {
	for _, e := range m.Entities {
		if len(e.Labels) == 0 {
			return errors.New("entities need at least 1 label specified")
		}
		if e.ExternalId == nil {
			return errors.New("entities need an external id")
		}
		if ContainsLabel(e.Labels, "DigitalTwin") {
			if e.TenantId == "" {
				return errors.New("digital twins need a tenant id specified")
			}
		}
	}
	return nil
}
