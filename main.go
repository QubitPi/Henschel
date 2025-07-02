package main

import (
	"flag"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"log"
	"net/http"
)

func kongHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	packerConfigFile := hclwrite.NewEmptyFile()
	config := packerConfigFile.Body()

	resourceBlock := config.AppendNewBlock("source", []string{"amazon-ebs", "qubitpi"})
	resourceBody := resourceBlock.Body()

	resourceBody.SetAttributeValue("ami_name", cty.StringVal("my-kong-api-gateway"))
	resourceBody.SetAttributeValue("force_deregister", cty.StringVal("true"))

	buildBlock := config.AppendNewBlock("build", nil)
	buildBody := buildBlock.Body()
	sources := []string{"source.amazon-ebs.qubitpi"}
	buildBody.SetAttributeValue("sources", cty.ListVal([]cty.Value{cty.StringVal(sources[0])}))

	provisionerBlock := buildBody.AppendNewBlock("provisioner", []string{"qubitpi-kong-api-gateway-provisioner"})
	provisionerBody := provisionerBlock.Body()
	provisionerBody.SetAttributeValue("homeDir", cty.StringVal("/home/ubuntu"))

	err := ioutil.WriteFile("output.tf", packerConfigFile.Bytes(), 0644)
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
