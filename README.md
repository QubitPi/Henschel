Henschel
========

[![GitHub Workflow Status]][GitHub Workflow URL]
[![Apache License Badge]][Apache License, Version 2.0]

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

License
-------

The use and distribution terms for [Henschel]() are covered by the [Apache License, Version 2.0].

[Apache License Badge]: https://img.shields.io/badge/Apache%202.0-FE5D26.svg?style=for-the-badge&logo=Apache&logoColor=white
[Apache License, Version 2.0]: https://www.apache.org/licenses/LICENSE-2.0

[GitHub Workflow Status]: https://img.shields.io/github/actions/workflow/status/QubitPi/Henschel/ci-cd.yaml?branch=master&logo=github&style=for-the-badge&label=CI/CD&labelColor=2088FF
[GitHub Workflow URL]: https://github.com/QubitPi/Henschel/actions/workflows/ci-cd.yaml

[HashiCorp Packer]: https://packer.qubitpi.org/packer
[HashiCorp Terraform]: https://packer.qubitpi.org/terraform
