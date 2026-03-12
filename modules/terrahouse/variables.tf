variable "user_uuid" {
  description = "The UUID of the user"
  type        = string
  default     = "c9b2d5f8-7e61-4a20-8b6a-b8e2f19f0d93"

  validation {
    condition     = length(regex("^[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}$", var.user_uuid)) > 0
    error_message = "user_uuid must be a valid UUID (example: 3fa85f64-5717-4562-b3fc-2c963f66afa6)."
  }
}

variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string

  validation {
    condition = can(regex("^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$", var.bucket_name))
    error_message = "Bucket name must be 3-63 characters, lowercase letters, numbers, dots or hyphens, and must start and end with a letter or number."
  }
}

variable "index_file_path" {
  description = "Path to the index.html file for the S3 static website"
  type        = string

  validation {
    condition     = can(fileexists(var.index_file_path))
    error_message = "The file path provided does not exist. Please provide a valid path to index.html."
  }
}

variable "error_file_path" {
  description = "Path to the error.html file for the S3 static website"
  type        = string

  validation {
    condition     = can(fileexists(var.error_file_path))
    error_message = "The file path provided does not exist. Please provide a valid path to error.html."
  }
}

variable "content_version" {
  description = "Version number used for content updates"
  type        = number

  validation {
    condition     = var.content_version > 0 && floor(var.content_version) == var.content_version
    error_message = "content_version must be a positive integer."
  }
}

variable "assets_path" {
  description = "Path to assets folder"
  type        = string

  validation {
    condition     = length(fileset(var.assets_path, "*")) > 0
    error_message = "The provided assets path does not exist or is not a directory."
  }
}