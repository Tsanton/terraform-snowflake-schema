# **terraform-snowflake-schema**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Requirements](#requirements)
- [Providers](#providers)
- [Modules](#modules)
- [Resources](#resources)
- [Inputs](#inputs)
- [Outputs](#outputs)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_snowflake"></a> [snowflake](#requirement\_snowflake) | >=0.51.0, < 1.0.0 |
| <a name="requirement_time"></a> [time](#requirement\_time) | >=0.8.0, < 1.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_snowflake"></a> [snowflake](#provider\_snowflake) | >=0.51.0, < 1.0.0 |
| <a name="provider_time"></a> [time](#provider\_time) | >=0.8.0, < 1.0.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [snowflake_database_grant.schema_access_role_database_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/database_grant) | resource |
| [snowflake_role.role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role) | resource |
| [snowflake_role_grants.schema_access_role_r_rw_inheritance_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_grants) | resource |
| [snowflake_role_grants.schema_access_role_rw_rwc_inheritance_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_grants) | resource |
| [snowflake_role_grants.schema_access_role_rwc_db_sys_admin_inheritance_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_grants) | resource |
| [snowflake_role_ownership_grant.role_ownership](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_ownership_grant) | resource |
| [snowflake_schema.schema](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/schema) | resource |
| [snowflake_schema_grant.schema_access_role_schema_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/schema_grant) | resource |
| [snowflake_schema_grant.schema_owner](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/schema_grant) | resource |
| [time_sleep.schema_ownership_wait](https://registry.terraform.io/providers/hashicorp/time/latest/docs/resources/sleep) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_db_name"></a> [db\_name](#input\_db\_name) | The name of the database under which to create this schema | `string` | n/a | yes |
| <a name="input_db_sys_admin_role_name"></a> [db\_sys\_admin\_role\_name](#input\_db\_sys\_admin\_role\_name) | The name database sys admin | `string` | n/a | yes |
| <a name="input_db_user_admin_role_name"></a> [db\_user\_admin\_role\_name](#input\_db\_user\_admin\_role\_name) | The name database role admin | `string` | n/a | yes |
| <a name="input_schema_name"></a> [schema\_name](#input\_schema\_name) | The name of the new schema | `string` | n/a | yes |
| <a name="input_tags_map"></a> [tags\_map](#input\_tags\_map) | Key value pairs that will be jsonencoded and stored as a json-string as Snowflake comments | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_access_roles"></a> [access\_roles](#output\_access\_roles) | map k, v where key in [r, rw, rwc] and the value equals the full access role name |
| <a name="output_base_access_role_name"></a> [base\_access\_role\_name](#output\_base\_access\_role\_name) | Access role base name in format <DB\_NAME>\_<SCHEMA\_NAME>\_ACCESS\_ROLE |
| <a name="output_schema_name"></a> [schema\_name](#output\_schema\_name) | The name of the schema |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
