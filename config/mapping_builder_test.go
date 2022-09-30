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

package config_test

import (
	"github.com/indykite/jarvis-sdk-go/config"
	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Mapping builder", func() {
	It("Build mapping of digital twin and another entity", func() {
		mappingBuilder := config.NewMappingBuilder()
		var dts []*configpb.IngestMappingConfig_Entity
		var entities []*configpb.IngestMappingConfig_Entity
		dtBuilder := config.NewDigitalTwinBuilder()
		dtBuilder.ExternalID("fodselsnummer").
			Properties([]*configpb.IngestMappingConfig_Property{
				{
					SourceName: "kallenavn",
					MappedName: "nickname",
					IsRequired: false,
				},
			}).Relationships([]*configpb.IngestMappingConfig_Relationship{
			{
				ExternalId: "familienummer",
				Type:       "MEMBER_OF",
				Direction:  configpb.IngestMappingConfig_DIRECTION_OUTBOUND,
				MatchLabel: "Family",
			},
			{
				ExternalId: "mors_fodselsnummer",
				Type:       "MOTHER_OF",
				Direction:  configpb.IngestMappingConfig_DIRECTION_INBOUND,
				MatchLabel: "DigitalTwin",
			},
		})

		dts = append(dts, dtBuilder.Build())

		familyBuilder := config.NewEntityBuilder()
		familyBuilder.Labels([]string{"Family"}).ExternalID("familienummer")

		entities = append(entities, familyBuilder.Build())

		mappingBuilder.Name("DSF mapping").DisplayName("some cool display name").
			Description("Description").DigitalTwins(dts).
			Entities(entities)
		mapping := mappingBuilder.Mapping
		Expect(mapping.Name).To(Equal("DSF mapping"))
		Expect(mapping.DisplayName).To(Equal("some cool display name"))
		Expect(mapping.Description).To(Equal("Description"))
		Expect(len(mapping.Entities)).To(Equal(2))

		dtMapping := mapping.Entities[0]
		Expect(dtMapping.GetLabels()).To(Equal([]string{"DigitalTwin"}))
		Expect(dtMapping.GetExternalId()).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"SourceName": Equal("fodselsnummer"),
			"MappedName": Equal("ExternalId"),
			"IsRequired": Equal(true),
		})))
		Expect(dtMapping.GetRelationships()[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ExternalId": Equal("familienummer"),
			"Type":       Equal("MEMBER_OF"),
			"Direction":  Equal(configpb.IngestMappingConfig_DIRECTION_OUTBOUND),
			"MatchLabel": Equal("Family"),
		})))

		familyEntity := mapping.Entities[1]
		Expect(familyEntity.GetLabels()).To(Equal([]string{"Family"}))
		Expect(familyEntity.GetExternalId()).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"SourceName": Equal("familienummer"),
			"MappedName": Equal("ExternalId"),
			"IsRequired": Equal(true),
		})))
		Expect(familyEntity.GetRelationships()).To(BeEmpty())
		Expect(familyEntity.GetProperties()).To(BeEmpty())
	})
})
