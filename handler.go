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
	"log"
	"net/http"
)

// KongHandler handles HTTP requests for deploying the Kong API Gateway.
// It expects a POST request with a JSON payload, parses it, and then
// generates a Packer HCL configuration file based on the payload.
func KongHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload KongApiGatewayPayload
	if err := GetJSONPayload(r, &payload); err != nil {
		log.Printf("Error parsing payload: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := GeneratePackerConfigFile(payload); err != nil {
		log.Printf("Error generating Packer config file: %v", err)
		http.Error(w, "Failed to generate Packer configuration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Packer configuration file 'kong.pkr.hcl' generated successfully."))
}
