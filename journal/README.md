# Terra Town Journal

## Root Model Structure

Root model structure for the Terra Town Journal is as follows:

```
PROJECT_ROOT
|
├── variables.tf
├── main.tf
├── provider.tf
├── outputs.tf
├── terraform.tfvars
└── README.md
```

## Terraform Variables

In terraform cloud we can set two types of variables:
1. Environment Variables
2. Terraform Variables

## Loading Terraform Input Variables

We can set variable during terraform plan using `-var` flag or we can use `terraform.tfvars` file to load variables automatically.

## Terraform Import

If a resource already exists in a cloud provider (AWS, Azure, etc.), you can import it into Terraform state so Terraform can manage it.

### How to import an existing resource?

1. Define the bucket in your Terraform code

    **main.tf**
    ```
    resource "aws_s3_bucket" "website_bucket" {
        bucket = var.bucket_name

        tags = {
            UserUUID = var.user_uuid
        }
    }
    ```

    **variables.tf**
    ```
    variable "bucket_name" {
        description = "Name of the S3 bucket"
        type        = string

        validation {
            condition = can(regex("^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$", var.bucket_name))
            error_message = "Bucket name must be 3-63 characters, lowercase letters, numbers, dots or hyphens, and must start and end with a letter or number."
        }
    }
    ```

    **outputs.tf**
    ```
    output "bucket_name" {
        value = aws_s3_bucket.website_bucket.bucket
    }
    ```

    **terraform.tfvars**
    ```
    bucket_name = "wali-yar-khan-bucket"
    ```


2. Run the import command

    ```terraform import aws_s3_bucket.website_bucket wali-yar-khan-bucket```

    **Format:**

    ```terraform import <resource_type>.<resource_name> <real_resource_id>```
    
    **So here:**
    - aws_s3_bucket → resource type
    - website_bucket → Terraform resource name
    - wali-yar-khan-bucket → actual AWS bucket name

3. Check configuration differences

    ```terraform plan```

### Fix Manual Configuration

If someone deletes or modifies the resource manually through ClickOps. Terraform will automatically detect and bring it back to the expecgted state.