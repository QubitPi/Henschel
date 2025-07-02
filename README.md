Henschel
========

__Henschel__ is a Go webservice that exposes [HashiCorp Packer] and [HashiCorp Terraform] as Services, i.e.
__HashiCorp as a Service__ (__HaaS__) and makes it possible to deploy any infrastructures in the matter of a single HTTP
request.

At its very early stage of development, Henschel supports generating AWS AMI Packer config files of

- Kong API gateway

Development
-----------

### Run Tests

```console
make test
```

### Spin Up the Service

```console
make run
```

Example request:

```console
curl -v -X POST localhost:8080/deployKongApiGateway \
    -H "Content-Type: application/json" \
    -d '{"sslCertBase64":"YXNkZnNnaHRkeWhyZXJ3ZGZydGV3ZHNmZ3RoeTY0cmV3ZGZyZWd0cmV3d2ZyZw==", "sslCertKeyBase64":"MzI0NXRnZjk4dmJoIGNsO2VbNDM1MHRdzszNDM1b2l0cmo=", "kongApiGatewayDomain":"api.mycompany.com"}'
```

At this moment, a "kong.pkr.hcl" config file for Packing an Kong AWS AMI will be generated at the project root.

[HashiCorp Packer]: https://packer.qubitpi.org/packer
[HashiCorp Terraform]: https://packer.qubitpi.org/terraform
