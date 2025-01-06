/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package http_generator

import (
	"fmt"
	"testing"

	"github.com/terance-edmonds/wso2-apk-k8s-go-lib/config/constants"
	"github.com/terance-edmonds/wso2-apk-k8s-go-lib/config/types"
	"github.com/terance-edmonds/wso2-apk-k8s-go-lib/pkg/utils"

	"github.com/stretchr/testify/assert"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// TestGenerator test for Generator
func TestGenerator(t *testing.T) {
	apkConf := types.APKConf{
		Name:                   "EmployeeServiceAPI",
		Version:                "3.14",
		BasePath:               "/employees-info",
		Type:                   "REST",
		DefaultVersion:         false,
		SubscriptionValidation: false,
		EndpointConfigurations: &types.EndpointConfigurations{
			Production: &types.EndpointConfiguration{
				Endpoint: types.EndpointURL("http://employee-service:8080"),
			},
		},
		RateLimit: &types.RateLimit{
			Unit:            "Minute",
			RequestsPerUnit: 5,
		},
		Authentication: &[]types.AuthConfiguration{
			{
				AuthType: "APIKey",
				Enabled:  true,
			},
		},
		Operations: &[]types.Operation{
			{Target: "/employees", Verb: "GET", Secured: true, Scopes: []string{}},
			{Target: "/employee", Verb: "POST", Secured: true, Scopes: []string{}},
			{Target: "/employee/{employeeId}", Verb: "PUT", Secured: true, Scopes: []string{}},
			{Target: "/employee/{employeeId}", Verb: "DELETE", Secured: true, Scopes: []string{}},
		},
	}
	organization := types.Organization{
		Name: "wso2",
	}
	gatewayConfig := types.GatewayConfigurations{
		Name:         "wso2-apim",
		ListenerName: "wso2-apim-gateway",
		Hostname:     "wso2-apim",
	}

	gen := Generator()

	// Get the endpoint to use
	endpoints := utils.GetEndpoints(apkConf)
	// If endpoints has production type
	if endpoint, ok := endpoints[constants.PRODUCTION_TYPE]; ok {
		httpRoute, err := gen.GenerateHTTPRoute(apkConf, organization, gatewayConfig, *apkConf.Operations, &endpoint, constants.PRODUCTION_TYPE, "unique-id", 1)
		if err != nil {
			fmt.Println(err)
		}

		assert.Nil(t, err)
		assert.IsType(t, &gwapiv1.HTTPRoute{}, httpRoute)
	}
}
