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
  endpoint = "http://localhost:4567"
  user_uuid = "e328f4ab-b99f-421c-84c9-4ccea042c7d1"
  token = "9b49b3fb-b8e9-483c-b703-97ba88eef8e0"
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