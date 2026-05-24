# TerraTown

**End-to-end Infrastructure-as-Code platform combining AWS cloud provisioning with a custom Terraform provider.**

[![Terraform](https://img.shields.io/badge/Terraform-1.x-7B42BC?style=flat-square&logo=terraform&logoColor=white)](https://terraform.io)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![AWS](https://img.shields.io/badge/AWS-S3%20%2B%20CloudFront-FF9900?style=flat-square&logo=amazonaws&logoColor=white)](https://aws.amazon.com)

---

## Overview

TerraTown is a full IaC workflow built in two parallel tracks:

1. **AWS static site infrastructure** — a reusable, variable-driven Terraform module that provisions S3, CloudFront, access control, and content lifecycle management
2. **Custom Terraform provider** — a Go-based `terratowns` provider that enables Terraform-managed resource operations against the TerraTowns API, integrated via a local provider mirror

Together they demonstrate practical cloud provisioning, modular Terraform design, custom provider development, and automated local tooling in one cohesive workflow.

---

## Features

### AWS infrastructure module
- S3 bucket provisioning configured for static website hosting
- CloudFront distribution setup for global content delivery
- Bucket access control and origin access policies
- Automated content lifecycle rules for object management
- Fully variable-driven — reusable across environments without modification

### Custom `terratowns` Terraform provider (Go)
- Implemented the full Terraform provider interface in Go
- Authenticated HTTP requests against the TerraTowns API
- Schema-based resource definition with create, read, update, and delete operations
- Clean Terraform state handling across all resource lifecycle events
- Integrated via local provider mirror with an automated build and install script

### Input validation
- Strong variable validation across the Terraform configuration
- Catches misconfigured inputs at `terraform plan` time, before any infrastructure is touched

---

## Project structure

```
.
├── modules/
│   └── terrahouse_aws/         # Reusable AWS static site module
│       ├── main.tf             # S3, CloudFront, lifecycle rules
│       ├── variables.tf        # Input definitions with validation
│       └── outputs.tf
├── terraform-provider-terratowns/   # Custom Go provider
│   ├── main.go
│   ├── resource_terratowns_home.go
│   └── go.mod
├── bin/
│   └── build_provider          # Automated build + local mirror install script
├── main.tf                     # Root config — calls module, uses custom provider
├── variables.tf
├── outputs.tf
└── terraform.tfvars.example
```

---

## Getting started

### Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.x
- [Go](https://go.dev/dl/) 1.21+
- AWS credentials configured (`aws configure` or environment variables)
- TerraTowns API credentials

### 1. Clone the repo

```bash
git clone https://github.com/YarKhan02/<repo-name>.git
cd <repo-name>
```

### 2. Build and install the custom provider

The build script compiles the Go provider and installs it into the local provider mirror automatically:

```bash
./bin/build_provider
```

### 3. Configure variables

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:

```hcl
# AWS
bucket_name     = "your-unique-bucket-name"
aws_region      = "us-east-1"
content_version = 1

# TerraTowns
terratowns_endpoint = "https://terratowns.cloud/api"
terratowns_token    = "your-api-token"
teacherseat_user_uuid = "your-uuid"
```

### 4. Deploy

```bash
terraform init
terraform plan
terraform apply
```

Terraform will provision the S3 bucket, CloudFront distribution, and TerraTowns resources in a single apply.

### 5. Destroy

```bash
terraform destroy
```

---

## Module usage

The `terrahouse_aws` module can be imported independently into any Terraform project:

```hcl
module "terrahouse_aws" {
  source = "./modules/terrahouse_aws"

  bucket_name     = var.bucket_name
  index_html_filepath = var.index_html_filepath
  error_html_filepath = var.error_html_filepath
  content_version = var.content_version
}
```

All variables include validation rules — Terraform will surface configuration errors at plan time before any AWS API calls are made.

---

## Custom provider: how it works

The `terratowns` provider is written in Go and registered through a [local filesystem mirror](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror) so Terraform resolves it without hitting the public registry.

```hcl
terraform {
  required_providers {
    terratowns = {
      source  = "local.providers/local/terratowns"
      version = "1.0.0"
    }
  }
}
```

The `bin/build_provider` script handles compilation and mirror placement in one step — no manual Go toolchain steps required after initial setup.

---

## Key concepts demonstrated

| Concept | Implementation |
|---|---|
| Modular Terraform design | Reusable `terrahouse_aws` module with variable-driven config |
| Custom provider development | Full Go provider with CRUD lifecycle against TerraTowns API |
| Input validation | `validation` blocks on all critical variables |
| Local provider mirror | Filesystem mirror + automated build script |
| Static site delivery | S3 + CloudFront + access control + lifecycle rules |
| State management | Clean Terraform state across custom provider resources |

---

## License

MIT License — see [LICENSE](LICENSE) for details.
