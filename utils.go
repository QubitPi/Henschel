// Copyright 2025 Jiaqi Liu. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetJSONPayload parses the JSON payload from an http.Request's body
// and decodes it into the provided interface `v`.
// It also performs validation if the `v` interface implements a Validate() error method.
func GetJSONPayload(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		if err == io.EOF {
			return fmt.Errorf("empty request body")
		}
		return fmt.Errorf("failed to decode JSON payload: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("request body contains unexpected extra data after JSON payload")
	}

	if validator, ok := v.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("missing or invalid required field: %w", err)
		}
	}

	return nil
}
