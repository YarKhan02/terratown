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

data "aws_caller_identity" "current" {}