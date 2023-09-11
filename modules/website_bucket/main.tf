data "aws_canonical_user_id" "current" {}

resource "aws_s3_bucket" "this" {
  bucket = var.name
}

resource "aws_s3_bucket_ownership_controls" "this" {
  bucket = aws_s3_bucket.this.id
  rule {
    object_ownership = "ObjectWriter"
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

locals {
  canonical_user_permissions = [
    "WRITE_ACP",
    "READ_ACP",
    "READ",
    "WRITE"
  ]

  all_user_permissions = [
    "READ",
    "READ_ACP",
    "WRITE"
  ]
}

resource "aws_s3_bucket_acl" "this" {
  depends_on = [aws_s3_bucket_ownership_controls.this, aws_s3_bucket_public_access_block.this]

  bucket = aws_s3_bucket.this.id

  access_control_policy {
    dynamic "grant" {
      for_each = local.canonical_user_permissions
      content {
        grantee {
          id   = data.aws_canonical_user_id.current.id
          type = "CanonicalUser"
        }
        permission = grant.value
      }
    }

    dynamic "grant" {
      for_each = local.all_user_permissions
      content {
        grantee {
          uri  = "http://acs.amazonaws.com/groups/global/AllUsers"
          type = "Group"
        }
        permission = grant.value
      }
    }

    owner {
      id = data.aws_canonical_user_id.current.id
    }
  }
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.id
  versioning_configuration {
    status = "Enabled"
  }
}
