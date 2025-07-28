// Copyright 2025 Jiaqi Liu. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
)

// KongApiGatewayPayload represents the structure of the expected JSON payload
// for configuring the Kong API Gateway.
// Fields are tagged with `json:"fieldName"` to correctly map JSON keys to struct fields.
type KongApiGatewayPayload struct {
	SslCertBase64        string `json:"sslCertBase64"`        // Base64 encoded SSL certificate
	SslCertKeyBase64     string `json:"sslCertKeyBase64"`     // Base64 encoded SSL certificate key
	KongApiGatewayDomain string `json:"kongApiGatewayDomain"` // Domain name for the Kong API Gateway
}

// Validate checks if all required fields of the KongApiGatewayPayload struct
// are present and have valid non-empty values.
// It returns an error if any required field is missing or empty.
func (u *KongApiGatewayPayload) Validate() error {
	if u.SslCertBase64 == "" {
		return fmt.Errorf("sslCertBase64 is a required field and cannot be empty")
	}
	if u.SslCertKeyBase64 == "" {
		return fmt.Errorf("sslCertKeyBase64 is a required field and cannot be empty")
	}
	if u.KongApiGatewayDomain == "" {
		return fmt.Errorf("kongApiGatewayDomain is a required field and cannot be empty")
	}
	return nil
}
