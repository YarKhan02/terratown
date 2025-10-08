terraform {
#   cloud { 
    
#     organization = "CherryTroll" 

#     workspaces { 
#       name = "terra-house" 
#     } 
#   } 

  required_providers {
    random = {
      source = "hashicorp/random"
      version = "3.7.2"
    }
    aws = {
      source = "hashicorp/aws"
      version = "6.15.0"
    }
  }
}

provider "aws" {
  region  = "ap-southeast-2"
  profile = "terra"
}
