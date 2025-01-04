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

package utils

import (
	"gw_artifacts/config/constants"
	"gw_artifacts/config/types"
	"regexp"
	"strconv"
	"strings"

	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// GetHost extracts the host from a given URL
func GetHost(endpoint types.Endpoint) string {
	var url string
	// check if the endpoint is a string or a K8sService
	switch v := endpoint.(type) {
	case types.EndpointURL:
		url = string(v)
	case types.K8sService:
		url = ConstructURlFromK8sService(v)
	}

	var host string
	if len(url) >= 8 && url[:8] == "https://" {
		host = url[8:]
	} else if len(url) >= 7 && url[:7] == "http://" {
		host = url[7:]
	} else {
		return ""
	}

	if indexOfColon := strings.Index(host, ":"); indexOfColon != -1 {
		return host[:indexOfColon]
	} else if indexOfSlash := strings.Index(host, "/"); indexOfSlash != -1 {
		return host[:indexOfSlash]
	} else {
		return host
	}
}

// GetPort extracts the port from a given URL
func GetPort(endpoint interface{}) int {
	var url string
	switch v := endpoint.(type) {
	case string:
		url = v
	default:
		url = ConstructURlFromK8sService(endpoint)
	}
	var hostPort string
	var protocol string
	if strings.HasPrefix(url, "https://") {
		hostPort = url[8:]
		protocol = "https"
	} else if strings.HasPrefix(url, "http://") {
		hostPort = url[7:]
		protocol = "http"
	} else {
		return -1
	}

	if indexOfSlash := strings.Index(hostPort, "/"); indexOfSlash != -1 {
		hostPort = hostPort[:indexOfSlash]
	}

	if indexOfColon := strings.Index(hostPort, ":"); indexOfColon != -1 {
		port := hostPort[indexOfColon+1:]
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return -1
		}
		return portInt
	} else {
		if protocol == "https" {
			return 443
		} else {
			return 80
		}
	}
}

// ConstructURlFromK8sService prepares the service URL from K8sService
func ConstructURlFromK8sService(endpoint interface{}) string {
	if k8sService, ok := endpoint.(types.K8sService); ok {
		return k8sService.Protocol + "://" + k8sService.Name + "." + k8sService.Namespace + ".svc.cluster.local:" + k8sService.Port
	} else {
		return ""
	}
}

// GetProtocol extracts the protocol from a given URL
func GetProtocol(endpoint interface{}) string {
	if k8sService, ok := endpoint.(types.K8sService); ok {
		if k8sService.Protocol == "" {
			return "http"
		}
		return k8sService.Protocol
	} else if strEndpoint, ok := endpoint.(string); ok {
		if strings.HasPrefix(strEndpoint, "https://") {
			return "https"
		} else {
			return "http"
		}
	}
	return "http"
}

// GetPath extracts the path from a given URL
func GetPath(url string) string {
	var hostPort string
	if strings.HasPrefix(url, "https://") {
		hostPort = url[8:]
	} else if strings.HasPrefix(url, "http://") {
		hostPort = url[7:]
	} else {
		return ""
	}

	if indexOfSlash := strings.Index(hostPort, "/"); indexOfSlash != -1 {
		return hostPort[indexOfSlash:]
	} else {
		return ""
	}
}

// RetrievePathPrefix generates a path prefix based on the operation and basePath
func RetrievePathPrefix(operation string, basePath string) string {
	splitValues := strings.Split(operation, "/")
	generatedPath := ""

	if operation == "/*" {
		return "(.*)"
	} else if operation == "/" {
		return "/"
	}

	re := regexp.MustCompile(`\{.*\}`)
	for _, pathPart := range splitValues {
		trimmedPathPart := strings.TrimSpace(pathPart)
		if len(trimmedPathPart) > 0 {
			// Path contains path param
			if re.MatchString(trimmedPathPart) {
				generatedPath += "/" + re.ReplaceAllString(trimmedPathPart, "(.*)")
			} else {
				generatedPath += "/" + trimmedPathPart
			}
		}
	}

	if strings.HasSuffix(generatedPath, "/*") {
		lastSlashIndex := strings.LastIndex(generatedPath[:len(generatedPath)-1], "/")
		generatedPath = generatedPath[:lastSlashIndex] + "(.*)"
	}

	return basePath + strings.TrimSpace(generatedPath)
}

// GeneratePrefixMatch generates a prefix match based on the endpoint and operation
func GeneratePrefixMatch(endpointToUse types.EndpointDetails, operation types.Operation) string {
	target := operation.Target
	if target == "" {
		target = "/*"
	}
	splitValues := strings.Split(target, "/")
	generatedPath := ""
	pathParamCount := 1

	if target == "/*" {
		generatedPath = "\\1"
	} else if target == "/" {
		generatedPath = "/"
	} else {
		for _, value := range splitValues {
			trimmedValue := strings.TrimSpace(value)
			if len(trimmedValue) > 0 {
				if matched, _ := regexp.MatchString("\\{.*\\}", trimmedValue); matched {
					generatedPath += "/" + regexp.MustCompile("\\{.*\\}").ReplaceAllString(trimmedValue, "\\"+strconv.Itoa(pathParamCount))
					pathParamCount++
				} else {
					generatedPath += "/" + trimmedValue
				}
			}
		}
	}

	if strings.HasSuffix(generatedPath, "/*") {
		lastSlashIndex := strings.LastIndex(generatedPath, "/")
		generatedPath = generatedPath[:lastSlashIndex] + "///" + strconv.Itoa(pathParamCount)
	}
	if endpointToUse.ServiceEntry {
		return strings.TrimSpace(generatedPath)
	}
	return generatedPath
}

func GetHostNames(apkConf types.APKConf, endpointType string, organization types.Organization) []gwapiv1.Hostname {
	// todo: need to implement this function
	var hosts []gwapiv1.Hostname
	environment := apkConf.Environment
	orgAndEnv := ""
	if environment != "" {
		orgAndEnv = orgAndEnv + "-" + environment
	}
	return hosts
}

// GetEndpoints retrieves the endpoint details from the provided APK configuration.
func GetEndpoints(apkConf types.APKConf) map[string]types.EndpointDetails {
	createdEndpoints := make(map[string]types.EndpointDetails)
	endpointConfigs := apkConf.EndpointConfigurations
	if endpointConfigs != nil {
		createdEndpoints = createEndpoints(endpointConfigs, "")
	}
	return createdEndpoints
}

// GetEndpointToUse returns the endpoint details based on the endpoint configurations and type.
func GetEndpointToUse(endpointConfigs *types.EndpointConfigurations, endpointType string) *types.EndpointDetails {
	if endpointConfigs != nil {
		operationLevelEndpoint := createEndpoints(endpointConfigs, endpointType)
		if _, ok := operationLevelEndpoint[endpointType]; ok {
			endpoint := operationLevelEndpoint[endpointType]
			return &endpoint
		}
	}
	return nil
}

// createEndpoints creates a map of endpoint details based on the provided configurations and endpoint type.
func createEndpoints(endpointConfigs *types.EndpointConfigurations, endpointType string) map[string]types.EndpointDetails {
	createdEndpoints := make(map[string]types.EndpointDetails)
	productionEndpointConfig := endpointConfigs.Production
	sandboxEndpointConfig := endpointConfigs.Sandbox
	if endpointType == constants.PRODUCTION_TYPE || productionEndpointConfig != nil {
		createdEndpoints[constants.PRODUCTION_TYPE] = types.EndpointDetails{
			Name: GetHost(productionEndpointConfig.Endpoint),
			URL:  ConstructURlFromK8sService(productionEndpointConfig.Endpoint),
		}
	}
	if endpointType == constants.SANDBOX_TYPE || sandboxEndpointConfig != nil {
		createdEndpoints[constants.SANDBOX_TYPE] = types.EndpointDetails{
			Name: GetHost(sandboxEndpointConfig.Endpoint),
			URL:  ConstructURlFromK8sService(sandboxEndpointConfig.Endpoint),
		}
	}
	return createdEndpoints
}
