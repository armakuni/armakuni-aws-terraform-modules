#!/bin/bash 

# script is to generate the terraform docs for our folder tree which is unsupported
# each of our modules live in modules/{module name}/src
# ├── README.md
# └── modules
#     ├── github_actions_oidc
#     ├── route53_zone
#     └── website_bucket
#         ├── README.md
#         ├── examples
#          │   └── website_bucket
#         ├── src
#         │   ├── main.tf
#         │   ├── outputs.tf
#         │   └── variables.tf
#         └── test
#             └── website_bucket
#                 └── website_bucket_test.go 

for module in modules/*/src; do
  module_name=$(basename $(dirname $module))
  echo "generating tf docs for module: $module_name"
  terraform-docs $module
done
