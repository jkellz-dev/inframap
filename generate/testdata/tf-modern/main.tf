# Uses Terraform 1.2+ language features that the old terraform v0.15 config
# loader rejected. inframap must still extract the resources and their links.
terraform {
  required_version = ">= 1.5"
}

variable "name" {
  type = string
  validation {
    condition     = length(var.name) > 0
    error_message = "name must not be empty"
  }
}

resource "aws_s3_bucket" "logs" {
  bucket = var.name

  lifecycle {
    precondition {
      condition     = var.name != "x"
      error_message = "Invalid name."
    }
    postcondition {
      condition     = self.bucket != ""
      error_message = "Bucket must have a name."
    }
  }
}

resource "aws_s3_bucket_policy" "logs" {
  bucket = aws_s3_bucket.logs.id
  policy = "{}"
}

import {
  to = aws_s3_bucket.logs
  id = "logs"
}

moved {
  from = aws_s3_bucket.old
  to   = aws_s3_bucket.logs
}

removed {
  from = aws_s3_bucket.gone
  lifecycle {
    destroy = false
  }
}

check "bucket" {
  assert {
    condition     = aws_s3_bucket.logs.bucket != ""
    error_message = "Bucket must have a name."
  }
}
