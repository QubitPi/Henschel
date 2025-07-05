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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetJSONPayload(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expected    *KongApiGatewayPayload
		expectError bool
	}{
		{
			name:    "Valid JSON payload",
			payload: `{"sslCertBase64": "foo", "sslCertKeyBase64": "bar", "kongApiGatewayDomain": "bat"}`,
			expected: &KongApiGatewayPayload{
				SslCertBase64:        "foo",
				SslCertKeyBase64:     "bat",
				KongApiGatewayDomain: "bat",
			},
			expectError: false,
		},
		{
			name:        "Invalid JSON payload (Malformed JSON)",
			payload:     `{"name": "Jane Doe", "email": "jane.doe@example.com",`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty JSON payload",
			payload:     `{}`,
			expected:    &KongApiGatewayPayload{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/deployKongApiGateway", bytes.NewBufferString(tt.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			var payload KongApiGatewayPayload
			parseErr := GetJSONPayload(req, &payload)

			if tt.expectError {
				if parseErr == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if parseErr != nil {
					t.Errorf("Did not expect an error, but got: %v", parseErr)
				}
				if &payload == nil {
					t.Fatalf("Expected a payload, but got nil")
				}
				if reflect.DeepEqual(payload, tt.expected) {
					t.Errorf("Expected payload %v, but got %v", tt.expected, payload)
				}
			}
		})
	}
}

func TestKongHandler(t *testing.T) {
	payload := make(map[string]string)
	payload["sslCertBase64"] = "YXNkZnNnaHRkeWhyZXJ3ZGZydGV3ZHNmZ3RoeTY0cmV3ZGZyZWd0cmV3d2ZyZw=="
	payload["sslCertKeyBase64"] = "MzI0NXRnZjk4dmJoIGNsO2VbNDM1MHRdzszNDM1b2l0cmo="
	payload["kongApiGatewayDomain"] = "api.mycompany.com"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/deployKongApiGateway", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	kongHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
