package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
