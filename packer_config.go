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
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GeneratePackerConfigFile creates a Packer HCL configuration file named "kong.pkr.hcl"
// based on the provided KongApiGatewayPayload.
// It orchestrates the creation of different HCL blocks by calling helper functions.
func GeneratePackerConfigFile(payload KongApiGatewayPayload) error {
	packerConfigFile := hclwrite.NewEmptyFile()
	config := packerConfigFile.Body()

	addRequiredPlugins(config)
	config.AppendNewline()

	addSourceBlock(config)
	config.AppendNewline()

	addBuildBlock(config, payload)

	err := os.WriteFile("kong.pkr.hcl", packerConfigFile.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write Packer configuration file: %v", err)
		return err
	}
	return nil
}

// addRequiredPlugins adds the 'required_plugins' block to the HCL body.
// This block defines the necessary Packer plugins (amazon and qubitpi).
func addRequiredPlugins(body *hclwrite.Body) {
	packerBlock := body.AppendNewBlock("required_plugins", nil)
	packerBody := packerBlock.Body()

	amazonBlock := packerBody.AppendNewBlock("amazon", nil)
	amazonBody := amazonBlock.Body()
	amazonBody.SetAttributeValue("version", cty.StringVal(">= 0.0.2"))
	amazonBody.SetAttributeValue("source", cty.StringVal("github.com/hashicorp/amazon"))

	qubitPiBlock := packerBody.AppendNewBlock("qubitpi", nil)
	qubitPiBody := qubitPiBlock.Body()
	qubitPiBody.SetAttributeValue("version", cty.StringVal(">= 0.0.50"))
	qubitPiBody.SetAttributeValue("source", cty.StringVal("github.com/QubitPi/qubitpi"))
}

// addSourceBlock adds the 'source' block to the HCL body.
// This block defines the source AMI and instance configuration for building the image.
func addSourceBlock(body *hclwrite.Body) {
	resourceBlock := body.AppendNewBlock("source", []string{"amazon-ebs", "qubitpi"})
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
}

// addBuildBlock adds the 'build' block to the HCL body.
// This block defines the build process and includes the 'provisioner' block.
func addBuildBlock(body *hclwrite.Body, payload KongApiGatewayPayload) {
	buildBlock := body.AppendNewBlock("build", nil)
	buildBody := buildBlock.Body()
	sources := []string{"source.amazon-ebs.qubitpi"} // Reference the defined source
	buildBody.SetAttributeValue("sources", cty.ListVal([]cty.Value{cty.StringVal(sources[0])}))

	provisionerBlock := buildBody.AppendNewBlock("provisioner", []string{"qubitpi-kong-api-gateway-provisioner"})
	provisionerBody := provisionerBlock.Body()
	provisionerBody.SetAttributeValue("homeDir", cty.StringVal("/home/ubuntu"))
	provisionerBody.SetAttributeValue("sslCertBase64", cty.StringVal(payload.SslCertBase64))
	provisionerBody.SetAttributeValue("sslCertKeyBase64", cty.StringVal(payload.SslCertKeyBase64))
	provisionerBody.SetAttributeValue("kongApiGatewayDomain", cty.StringVal(payload.KongApiGatewayDomain))
}
