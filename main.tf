# terraform {
#   cloud { 
    
#     organization = "CherryTroll" 

#     workspaces { 
#       name = "terra-house" 
#     } 
#   } 
# }

module "terrahouse" {
  source = "./modules/terrahouse"
  user_uuid = var.user_uuid
  bucket_name = var.bucket_name
}