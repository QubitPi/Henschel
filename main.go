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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// KongApiGatewayPayload represents the structure of your expected JSON payload.
// Adjust fields and types according to your actual JSON structure.
type KongApiGatewayPayload struct {
	SslCertBase64        string `json:"sslCertBase64"`
	SslCertKeyBase64     string `json:"sslCertKeyBase64"`
	KongApiGatewayDomain string `json:"kongApiGatewayDomain"`
}

// GetJSONPayload parses the JSON payload from an http.Request and returns it.
// It decodes the JSON into the provided `v` interface.
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

	return nil
}

func kongHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload KongApiGatewayPayload
	err := GetJSONPayload(r, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	packerConfigFile := hclwrite.NewEmptyFile()
	config := packerConfigFile.Body()

	packerBlock := config.AppendNewBlock("required_plugins", nil)
	packerBody := packerBlock.Body()
	amazonBlock := packerBody.AppendNewBlock("amazon", nil)
	amazonBody := amazonBlock.Body()
	amazonBody.SetAttributeValue("version", cty.StringVal(">= 0.0.2"))
	amazonBody.SetAttributeValue("source", cty.StringVal("github.com/hashicorp/amazon"))
	qubitPiBlock := packerBody.AppendNewBlock("qubitpi", nil)
	qubitPiBody := qubitPiBlock.Body()
	qubitPiBody.SetAttributeValue("version", cty.StringVal(">= 0.0.50"))
	qubitPiBody.SetAttributeValue("source", cty.StringVal("github.com/QubitPi/qubitpi"))

	config.AppendNewline()

	resourceBlock := config.AppendNewBlock("source", []string{"amazon-ebs", "qubitpi"})
	resourceBody := resourceBlock.Body()
	resourceBody.SetAttributeValue("ami_name", cty.StringVal("my-kong-api-gateway"))
	resourceBody.SetAttributeValue("force_deregister", cty.StringVal("true"))
	resourceBody.SetAttributeValue("force_delete_snapshot", cty.StringVal("true"))
	resourceBody.SetAttributeValue("instance_type", cty.StringVal("t2.micro"))
	resourceBody.SetAttributeValue("region", cty.StringVal("us-west-1"))
	resourceBody.SetAttributeValue("ssh_username", cty.StringVal("ubuntu"))
	sourceAmiFilterBlock := resourceBody.AppendNewBlock("source_ami_filter", nil)
	sourceAmiFilterBody := sourceAmiFilterBlock.Body()
	sourceAmiFilterBody.SetAttributeValue("most_recent", cty.BoolVal(true))
	sourceAmiFilterBody.SetAttributeValue("owners", cty.ListVal([]cty.Value{cty.StringVal("099720109477")}))
	filtersBlock := sourceAmiFilterBody.AppendNewBlock("filters", nil)
	filtersBody := filtersBlock.Body()
	filtersBody.SetAttributeValue("name", cty.StringVal("ubuntu/images/*ubuntu-*-22.04-amd64-server-*"))
	filtersBody.SetAttributeValue("root-device-type", cty.StringVal("ebs"))
	filtersBody.SetAttributeValue("virtualization-type", cty.StringVal("hvm"))
	launchBlockDeviceMappingsBlock := resourceBody.AppendNewBlock("launch_block_device_mappings", nil)
	launchBlockDeviceMappingsBody := launchBlockDeviceMappingsBlock.Body()
	launchBlockDeviceMappingsBody.SetAttributeValue("device_name", cty.StringVal("/dev/sda1"))
	launchBlockDeviceMappingsBody.SetAttributeValue("volume_size", cty.NumberIntVal(8))
	launchBlockDeviceMappingsBody.SetAttributeValue("volume_type", cty.StringVal("gp2"))
	launchBlockDeviceMappingsBody.SetAttributeValue("delete_on_termination", cty.BoolVal(true))

	config.AppendNewline()

	buildBlock := config.AppendNewBlock("build", nil)
	buildBody := buildBlock.Body()
	sources := []string{"source.amazon-ebs.qubitpi"}
	buildBody.SetAttributeValue("sources", cty.ListVal([]cty.Value{cty.StringVal(sources[0])}))

	provisionerBlock := buildBody.AppendNewBlock("provisioner", []string{"qubitpi-kong-api-gateway-provisioner"})
	provisionerBody := provisionerBlock.Body()
	provisionerBody.SetAttributeValue("homeDir", cty.StringVal("/home/ubuntu"))
	provisionerBody.SetAttributeValue("sslCertBase64", cty.StringVal(payload.SslCertBase64))
	provisionerBody.SetAttributeValue("sslCertKeyBase64", cty.StringVal(payload.SslCertKeyBase64))
	provisionerBody.SetAttributeValue("kongApiGatewayDomain", cty.StringVal(payload.KongApiGatewayDomain))

	err = ioutil.WriteFile("kong.pkr.hcl", packerConfigFile.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.String("port", ":8080", "Webservice port; default to 8080")
	flag.Parse()

	http.HandleFunc("/deployKongApiGateway", kongHandler)

	http.ListenAndServe(*port, nil)
}
