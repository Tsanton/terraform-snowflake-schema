locals {
  role_suffixes = ["r", "rw", "rwc"]
}

resource "snowflake_role" "role" {
  for_each = toset(local.role_suffixes)

  name    = upper("${local.base_access_role_name}_${each.key}")
  comment = upper(jsonencode(merge(var.tags_map, {})))

  depends_on = [
    snowflake_schema.schema
  ]

  lifecycle {
    replace_triggered_by = [
      snowflake_schema.schema.database,
      snowflake_schema.schema.name
    ]
  }
}

resource "snowflake_role_ownership_grant" "role_ownership" {
  for_each = snowflake_role.role

  on_role_name = each.value.name
  to_role_name = var.db_user_admin_role_name

  current_grants = "COPY"

  revert_ownership_to_role_name = var.db_sys_admin_role_name
}

resource "snowflake_role_grants" "schema_access_role_r_rw_inheritance_grant" {
  role_name = snowflake_role.role["r"].name

  roles = [
    snowflake_role.role["rw"].name
  ]

  users = []

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]

  lifecycle {
    replace_triggered_by = [
      snowflake_schema.schema.database,
      snowflake_schema.schema.name
    ]
  }
}

resource "snowflake_role_grants" "schema_access_role_rw_rwc_inheritance_grant" {
  role_name = snowflake_role.role["rw"].name

  roles = [
    snowflake_role.role["rwc"].name
  ]

  users = []

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]

  lifecycle {
    replace_triggered_by = [
      snowflake_schema.schema.database,
      snowflake_schema.schema.name
    ]
  }
}


resource "snowflake_role_grants" "schema_access_role_rwc_db_sys_admin_inheritance_grant" {
  role_name = snowflake_role.role["rwc"].name

  roles = [
    var.db_sys_admin_role_name
  ]

  users = []

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]

  lifecycle {
    replace_triggered_by = [
      snowflake_schema.schema.database,
      snowflake_schema.schema.name
    ]
  }
}

resource "snowflake_database_grant" "schema_access_role_database_grant" {
  database_name = var.db_name

  privilege = "USAGE"
  roles     = [for val in snowflake_role.role : val.name]

  with_grant_option      = false
  enable_multiple_grants = true

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]

  lifecycle {
    replace_triggered_by = [
      snowflake_schema.schema.database,
      snowflake_schema.schema.name
    ]
  }
}

resource "snowflake_schema_grant" "schema_access_role_schema_grant" {
  database_name = var.db_name
  schema_name   = snowflake_schema.schema.name

  privilege = "USAGE"
  roles     = [for val in snowflake_role.role : val.name]

  with_grant_option      = false
  enable_multiple_grants = true

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]
}

//sleep to avoid race condition issues
resource "time_sleep" "schema_ownership_wait" {
  create_duration = "10s"

  depends_on = [
    snowflake_role_ownership_grant.role_ownership
  ]
}
