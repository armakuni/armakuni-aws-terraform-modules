provider "aws" {
  region = var.region
}

# NOTE: S3 bucket created in region based on the aws provider set,
#       currently no way as of to directly configure in the s3 resource
module "test_website_bucket" {
  source = "../.."
  name   = var.name
}
