terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.53.0"
    }
  }
}
variable "snowflake_uid" { type = string }
variable "snowflake_pwd" { type = string }
variable "snowflake_account" { type = string }
variable "snowflake_region" { type = string }

provider "snowflake" {
  username = var.snowflake_uid
  password = var.snowflake_pwd
  account  = var.snowflake_account
  region   = var.snowflake_region
  role     = "ACCOUNTADMIN"
}

variable "db_prefix" { type = string }
variable "db_identifier" { type = string }

module "db" {
  source        = "github.com/tsanton/terraform-snowflake-database?ref=v1.0.3"
  db_prefix     = var.db_prefix
  db_identifier = var.db_identifier
}

variable "schema_name" { type = string }

module "schema" {
  source      = "../../"
  db_name     = module.db.db_name
  schema_name = var.schema_name

  db_user_admin_role_name = module.db.user_admin_name
  db_sys_admin_role_name  = module.db.sys_admin_name

  depends_on = [
    module.db
  ]
}

module "schema_static" {
  source      = "../../"
  db_name     = module.db.db_name
  schema_name = "STATIC"

  db_user_admin_role_name = module.db.user_admin_name
  db_sys_admin_role_name  = module.db.sys_admin_name

  depends_on = [
    module.db
  ]
}
