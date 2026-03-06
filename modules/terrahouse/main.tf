terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "6.34.0"
    }
  }
}

provider "aws" {
  region  = "us-east-1"
  profile = "terra"
}


resource "aws_s3_bucket" "website_bucket" {
  bucket = var.bucket_name

  tags = {
    UserUUID = var.user_uuid
  }
}

