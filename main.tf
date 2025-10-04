terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
      version = "3.7.2"
    }
  }
}

resource "random_string" "bucket_name" {
  length           = 16
  special          = false
}

output "random_bucket_name" {
    value = random_string.bucket_name.result
}