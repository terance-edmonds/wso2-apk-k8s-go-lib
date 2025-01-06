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
	"strconv"

	"github.com/terance-edmonds/wso2-apk-k8s-go-lib/config/types"
	"github.com/terance-edmonds/wso2-apk-k8s-go-lib/pkg/utils"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HttpRouteGenerator is the interface for the HTTP route generator.
type httpRouteGenerator struct {
	GenerateHTTPRouteRules        func(apkConf types.APKConf, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string) ([]gwapiv1.HTTPRouteRule, error)
	GenerateHTTPRouteRule         func(apkConf types.APKConf, operation types.Operation, endpoint *types.EndpointDetails, endpointType string) (*gwapiv1.HTTPRouteRule, error)
	GenerateAndRetrieveParentRefs func(gatewayConfig types.GatewayConfigurations, uniqueId string) []gwapiv1.ParentReference
	GenerateHTTPRouteFilters      func(apkConf types.APKConf, endpointToUse types.EndpointDetails, operation types.Operation, endpointType string) ([]gwapiv1.HTTPRouteFilter, bool)
	ExtractHTTPRouteFilter        func(apkConf *types.APKConf, endpoint types.EndpointDetails, operation types.Operation, operationPolicies []types.OperationPolicy, isRequest bool) ([]gwapiv1.HTTPRouteFilter, bool)
	GetHostNames                  func(apkConf types.APKConf, endpointType string, organization types.Organization) []gwapiv1.Hostname
	RetrieveHTTPMatches           func(apkConf types.APKConf, operation types.Operation) ([]gwapiv1.HTTPRouteMatch, error)
	RetrieveHTTPMatch             func(apkConf types.APKConf, operation types.Operation) (gwapiv1.HTTPRouteMatch, error)
	GenerateHTTPBackEndRef        func(endpoint types.EndpointDetails, operation types.Operation, endpointType string) []gwapiv1.HTTPBackendRef
}

// Generator creates a new HTTP route generator.
func Generator() *httpRouteGenerator {
	gen := &httpRouteGenerator{}
	gen.GenerateHTTPRouteRules = gen.generateHTTPRouteRules
	gen.GenerateHTTPRouteRule = gen.generateHTTPRouteRule
	gen.GenerateAndRetrieveParentRefs = gen.generateAndRetrieveParentRefs
	gen.GenerateHTTPRouteFilters = gen.generateHTTPRouteFilters
	gen.ExtractHTTPRouteFilter = gen.extractHTTPRouteFilter
	gen.GetHostNames = utils.GetHostNames
	gen.RetrieveHTTPMatches = gen.retrieveHTTPMatches
	gen.RetrieveHTTPMatch = gen.retrieveHTTPMatch
	gen.GenerateHTTPBackEndRef = gen.generateHTTPBackEndRef
	return gen
}

// GenerateHTTPRoute generates a HTTPRoute based on the provided configurations.
func (g *httpRouteGenerator) GenerateHTTPRoute(apkConf types.APKConf, organization types.Organization, gatewayConfiguration types.GatewayConfigurations, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string, uniqueId string, count int) (*gwapiv1.HTTPRoute, error) {
	httpRouteRules, err := g.GenerateHTTPRouteRules(apkConf, operations, endpoint, endpointType)
	if err != nil {
		return nil, err
	}
	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: v1.ObjectMeta{
			Name: uniqueId + "-" + endpointType + "-httproute-" + strconv.Itoa(count),
		},
		Spec: gwapiv1.HTTPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: g.GenerateAndRetrieveParentRefs(gatewayConfiguration, uniqueId),
			},
			Rules:     httpRouteRules,
			Hostnames: g.GetHostNames(apkConf, endpointType, organization),
		},
	}
	return &httpRoute, nil
}
