Henschel
========

__Henschel__ is a HashiCorp webservice that exposes [HashiCorp Packer]() and [HashiCorp Terraform]() as Services, i.e.
__HashiCorp as a Service__ (__HaaS__) and makes it possible to deploy any infrastructures in a matter of single HTTP
request.

Development
-----------

```console
make run
```

```console
curl -v -X POST localhost:8080/deployKongApiGateway
```

[HashiCorp Packer]: https://packer.qubitpi.org/packer
[HashiCorp Terraform]: https://packer.qubitpi.org/terraform
