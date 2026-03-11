output "bucket_name" {
    description = "Bucket name for static website"
    value = module.terrahouse.bucket_name
}

output "website_endpoint" {
    description = "Website endpoint for the S3 bucket"
    value = module.terrahouse.website_endpoint
}

output "cloudfront" {
    description = "CloudFront distribution domain name"
    value = module.terrahouse.cloudfront
}