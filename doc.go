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

//go:generate mockgen -copyright_file ./doc/LICENSE -package config -destination ./test/config/v1beta1/config_management_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1 ConfigManagementAPIClient,ConfigManagementAPI_ListApplicationSpacesClient,ConfigManagementAPI_ListApplicationsClient,ConfigManagementAPI_ListApplicationAgentsClient
//go:generate mockgen -copyright_file ./doc/LICENSE -package ingest -destination ./test/ingest/v1beta3/ingest_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3 IngestAPIClient,IngestAPI_StreamRecordsClient
//go:generate mockgen -copyright_file ./doc/LICENSE -package authorization -destination ./test/authorization/v1beta1/authorization_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1 AuthorizationAPIClient
//go:generate mockgen -copyright_file ./doc/LICENSE -package knowledge -destination ./test/knowledge/v1beta2/identity_knowledge_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2 IdentityKnowledgeAPIClient
//go:generate mockgen -copyright_file ./doc/LICENSE -package tda -destination ./test/tda/v1beta1/trusted_data_access_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1 TrustedDataAccessAPIClient
//go:generate mockgen -copyright_file ./doc/LICENSE -package entitymatching -destination ./test/entitymatching/v1beta1/entitymatching_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1 EntityMatchingAPIClient

/*
Package indykite is the root of the packages used to access IndyKite Platform.

Debugging
To see gRPC logs, set the environment variable GRPC_GO_LOG_SEVERITY_LEVEL. See
https://pkg.go.dev/google.golang.org/grpc/grpclog for more information.
For HTTP logging, set the GODEBUG environment variable to "http2debug=1" or "http2debug=2".
*/
package indykite
