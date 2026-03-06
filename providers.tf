terraform {
#   cloud { 
    
#     organization = "CherryTroll" 

#     workspaces { 
#       name = "terra-house" 
#     } 
#   } 

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
