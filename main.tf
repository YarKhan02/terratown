terraform {
  required_providers {
    terratowns = {
      source = "local.providers/local/terratowns"
      version = "1.0.0"
    }
  }
#   cloud { 
    
#     organization = "CherryTroll" 

#     workspaces { 
#       name = "terra-house" 
#     } 
#   } 
}

provider "terratowns" {
  endpoint = var.terratowns_endpoint
  user_uuid = var.user_uuid
  token = "509df6e8-7950-49c3-b980-9e041b958bfb"
}

module "terrahouse" {
  source = "./modules/terrahouse"
  user_uuid = var.user_uuid
  bucket_name = var.bucket_name
  index_file_path = var.index_file_path
  error_file_path = var.error_file_path
  content_version = var.content_version
  assets_path = var.assets_path
}

resource "terratowns_home" "home" {
  name = "How to play acranum in 2023"
  description = "Acranum is a game from 2001"
  domain_name = "3fdq3gz.cloudfront.net"
  town = "missingo"
  content_version = "1"
}