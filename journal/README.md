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
    ```sh
    resource "aws_s3_bucket" "website_bucket" {
        bucket = var.bucket_name

        tags = {
            UserUUID = var.user_uuid
        }
    }
    ```

    **variables.tf**
    ```sh
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
    ```sh
    output "bucket_name" {
        value = aws_s3_bucket.website_bucket.bucket
    }
    ```

    **terraform.tfvars**
    ```sh
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

## Module Structure

[Module Structure](https://developer.hashicorp.com/terraform/language/v1.15.x/modules/develop/structure)

```
├── README.md
├── main.tf
├── variables.tf
├── outputs.tf
├── ...
├── modules/
│   ├── nestedA/
│   │   ├── README.md
│   │   ├── variables.tf
│   │   ├── main.tf
│   │   ├── outputs.tf
│   ├── nestedB/
│   ├── .../
├── examples/
│   ├── exampleA/
│   │   ├── main.tf
│   ├── exampleB/
│   ├── .../
```

### Module Sources

[Module Sources](https://developer.hashicorp.com/terraform/language/v1.15.x/modules/configuration)

Using the resource we can import the module from various places e.g
- Local
- GitHub
- Terraform Registry

```sh
module "terrahouse" {
  source = "./modules/terrahouse"
  user_uuid = var.user_uuid
  bucket_name = var.bucket_name
}
```

### Variable Declaration

We also have to declare variables in the root `variables.tf` file. We don't need to declare them completely with validation, but just simply declare them with type and description. The validation will be done in the module itself.

## S3 Website Hosting

### S3 Bucket Website Configuration
[Website Configure](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration)
```sh
resource "aws_s3_bucket_website_configuration" "website_config" {
  bucket = aws_s3_bucket.website_bucket.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}
```

### S3 Bucket Object Upload
[Bucket Object](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_object)
```sh
resource "aws_s3_object" "object" {
  bucket = aws_s3_bucket.website_bucket.bucket
  key    = "new_object_key"
  source = "path/to/file"

  etag = filemd5("path/to/file")
}
```

**-> etag**

Terraform do not automatically update the object if the file content changes. To force update we can use `etag` argument with `filemd5` cryptographic function.

### Output Website Endpoint

```
output "website_endpoint" {
    value = aws_s3_bucket_website_configuration.website_config.website_endpoint
}
```

## Working with files

### File exists function

```sh
variable "index_file_path" {
  description = "Path to the index.html file for the S3 static website"
  type        = string

  validation {
    condition     = can(fileexists(var.index_file_path))
    error_message = "The file path provided does not exist. Please provide a valid path to index.html."
  }
}
```

### Path Variable

[Special Path Variable]([Website Configure](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration))

In terraform there is a special variable called `path` that allows to reference local paths
- path.module → path to the current module
- path.root → path to the root module