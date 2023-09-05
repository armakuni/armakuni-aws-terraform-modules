# Armakuni AWS Terraform Modules

Terraform IaC Terraform modules used throughout Armakuni's AWS 

> TLDR; See modules directory, for each module contains Terraform Documentation

## Modules

| Module                                              | Description                                                                   |
|-----------------------------------------------------|-------------------------------------------------------------------------------|
| [github_actions_oidc](modules/github_actions_oidc/) | Setup a Github OIDC assosiated role with config for GH org & repo             |
| [route53_zone](modules/route53_zone/)               | Setup a domain and dns with externalised configuration                        |
| [website_bucket](modules/website_bucket/)           | S3 Bucket with sensible configuration and ACL support for serving websites    |

## Testing

This repository leverages [Terratest](https://terratest.gruntwork.io/docs/getting-started/introduction/) and automates the testing process via [Github Actions pipeline](.github/workflows/terratest_module_website_bucket.yml).

Terratest strongly recommend running tests in an environment that is totally separate from production, with this recomendation we make use of the training environments.
