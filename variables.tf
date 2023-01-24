variable "db_name" {
  description = "The name of the database under which to create this schema"
  type        = string
}

variable "schema_name" {
  description = "The name of the new schema"
  type        = string
}

variable "tags_map" {
  description = "Key value pairs that will be jsonencoded and stored as a json-string as Snowflake comments"
  type        = map(string)
  default     = {}
}

variable "db_user_admin_role_name" {
  description = "The name database role admin"
  type        = string
}

variable "db_sys_admin_role_name" {
  description = "The name database sys admin"
  type        = string
}

locals {
  base_access_role_name = upper("${var.db_name}_${var.schema_name}_ACCESS_ROLE")
}
