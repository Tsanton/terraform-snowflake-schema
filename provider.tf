terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">=0.51.0, < 1.0.0"
    }
    time = {
      source  = "hashicorp/time"
      version = ">=0.8.0, < 1.0.0"
    }
  }
}
