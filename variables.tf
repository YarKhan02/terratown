variable "user_uuid" {
  description = "The UUID of the user"
  type        = string
  default     = "c9b2d5f8-7e61-4a20-8b6a-b8e2f19f0d93"

  validation {
    condition     = length(regex("^[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}$", var.user_uuid)) > 0
    error_message = "user_uuid must be a valid UUID (example: 3fa85f64-5717-4562-b3fc-2c963f66afa6)."
  }
}
