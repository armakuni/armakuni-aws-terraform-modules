# Armakuni AWS Terraform Modules

Terraform IaC Terraform modules used throughout Armakuni's AWS 

> TLDR; See modules directory, for each module contains Terraform Documentation

## Modules

| Module              | Description                                                                   |
|---------------------|-------------------------------------------------------------------------------|
| github_actions_oidc |                                                                               |
| route53_zone        |                                                                               |
| website_bucket      |                                                                               |

## Testing

This repository leverages [Terratest](https://terratest.gruntwork.io/docs/getting-started/introduction/) and automates the testing process via [Github Actions pipeline](.github/workflows/terraform_test.yml).

Terratest strongly recommend running tests in an environment that is totally separate from production, with this recomendation we make use of the training environments.

<!-- BEGIN_TF_DOCS -->

<!-- END_TF_DOCS -->