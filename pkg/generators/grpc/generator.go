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

package grpc_generator

import (
	"gw_artifacts/config/types"
	"gw_artifacts/pkg/utils"
	"strconv"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// grpcRouteGenerator is the interface for the GRPC route generator.
type grpcRouteGenerator struct {
	GenerateGRPCRouteRules        func(apkConf types.APKConf, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string) ([]gwapiv1.GRPCRouteRule, error)
	GenerateGRPCRouteRule         func(apkConf types.APKConf, operation types.Operation, endpoint *types.EndpointDetails, endpointType string) (*gwapiv1.GRPCRouteRule, error)
	GenerateAndRetrieveParentRefs func(gatewayConfig types.GatewayConfigurations, uniqueId string) []gwapiv1.ParentReference
	GetHostNames                  func(apkConf types.APKConf, endpointType string, organization types.Organization) []gwapiv1.Hostname
	RetrieveGRPCMatches           func(operation types.Operation) []gwapiv1.GRPCRouteMatch
	RetrieveGRPCMatch             func(operation types.Operation) gwapiv1.GRPCRouteMatch
	GenerateGRPCBackEndRef        func(endpoint types.EndpointDetails, operation types.Operation) []gwapiv1.GRPCBackendRef
}

// Generator creates a new GRPC route generator.
func Generator() *grpcRouteGenerator {
	gen := &grpcRouteGenerator{}
	gen.GenerateGRPCRouteRules = gen.generateGRPCRouteRules
	gen.GenerateGRPCRouteRule = gen.generateGRPCRouteRule
	gen.GenerateAndRetrieveParentRefs = gen.generateAndRetrieveParentRefs
	gen.GetHostNames = utils.GetHostNames
	gen.RetrieveGRPCMatches = gen.retrieveGRPCMatches
	gen.RetrieveGRPCMatch = gen.retrieveGRPCMatch
	gen.GenerateGRPCBackEndRef = gen.generateGRPCBackEndRef
	return gen
}

// GenerateGRPCRoute generates a GRPCRoute based on the provided configurations.
func (g *grpcRouteGenerator) GenerateGRPCRoute(apkConf types.APKConf, organization types.Organization, gatewayConfiguration types.GatewayConfigurations, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string, uniqueId string, count int) (*gwapiv1.GRPCRoute, error) {
	grpcRouteRules, err := g.GenerateGRPCRouteRules(apkConf, operations, endpoint, endpointType)
	if err != nil {
		return nil, err
	}
	grpcRoute := gwapiv1.GRPCRoute{
		ObjectMeta: v1.ObjectMeta{
			Name: uniqueId + "-" + endpointType + "-grpcroute-" + strconv.Itoa(count),
		},
		Spec: gwapiv1.GRPCRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: g.GenerateAndRetrieveParentRefs(gatewayConfiguration, uniqueId),
			},
			Rules:     grpcRouteRules,
			Hostnames: g.GetHostNames(apkConf, endpointType, organization),
		},
	}
	return &grpcRoute, nil
}
