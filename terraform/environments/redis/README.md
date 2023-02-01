<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.4.3 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_kind"></a> [kind](#module\_kind) | ../../modules/kind | n/a |
| <a name="module_mysql"></a> [mysql](#module\_mysql) | ../../modules/mysql | n/a |
| <a name="module_postgresql"></a> [postgresql](#module\_postgresql) | ../../modules/postgresql | n/a |
| <a name="module_redis"></a> [redis](#module\_redis) | ../../modules/redis | n/a |
| <a name="module_sql_mysql"></a> [sql\_mysql](#module\_sql\_mysql) | ../../modules/sql | n/a |
| <a name="module_sql_postgres"></a> [sql\_postgres](#module\_sql\_postgres) | ../../modules/sql | n/a |

## Resources

| Name | Type |
|------|------|
| [random_pet.mysql_database](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.mysql_user](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.postgres_database](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.postgres_user](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_mysql_database"></a> [mysql\_database](#output\_mysql\_database) | MySQL database |
| <a name="output_mysql_password"></a> [mysql\_password](#output\_mysql\_password) | MySQL password |
| <a name="output_mysql_url"></a> [mysql\_url](#output\_mysql\_url) | MySQL URL |
| <a name="output_postgres_database"></a> [postgres\_database](#output\_postgres\_database) | Postgres database |
| <a name="output_postgres_password"></a> [postgres\_password](#output\_postgres\_password) | Postgres password |
| <a name="output_postgresql_url"></a> [postgresql\_url](#output\_postgresql\_url) | Postgres URL |
| <a name="output_redis_password"></a> [redis\_password](#output\_redis\_password) | Redis password |
| <a name="output_redis_url"></a> [redis\_url](#output\_redis\_url) | Redis URL |
| <a name="output_sql_mysql_url"></a> [sql\_mysql\_url](#output\_sql\_mysql\_url) | SQL MySQL URL |
| <a name="output_sql_postgres_url"></a> [sql\_postgres\_url](#output\_sql\_postgres\_url) | SQL Postgres URL |
<!-- END_TF_DOCS -->