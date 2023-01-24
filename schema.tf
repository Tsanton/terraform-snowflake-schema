resource "snowflake_schema" "schema" {
  database = var.db_name
  name     = upper(var.schema_name)
  comment  = upper(jsonencode(merge(var.tags_map, {})))

  is_transient        = false
  is_managed          = true
  data_retention_days = 1
}

resource "snowflake_schema_grant" "schema_owner" {
  database_name = var.db_name
  schema_name   = snowflake_schema.schema.name

  privilege = "OWNERSHIP"
  roles = [
    snowflake_role.role["rwc"].name
  ]

  with_grant_option      = false
  enable_multiple_grants = true

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]
}
