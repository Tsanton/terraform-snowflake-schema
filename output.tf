output "base_access_role_name" {
  description = "Access role base name in format <DB_NAME>_<SCHEMA_NAME>_ACCESS_ROLE"
  value       = local.base_access_role_name
}

output "schema_name" {
  description = "The name of the schema"
  value       = snowflake_schema.schema.name
}

output "access_roles" {
  description = "map k, v where key in [r, rw, rwc] and the value equals the full access role name"
  value       = { for k, v in snowflake_role.role : k => v.name }
}
