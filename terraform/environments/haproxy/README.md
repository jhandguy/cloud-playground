<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.5.1 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_certmanager"></a> [certmanager](#module\_certmanager) | ../../modules/certmanager | n/a |
| <a name="module_haproxy"></a> [haproxy](#module\_haproxy) | ../../modules/haproxy | n/a |
| <a name="module_kind"></a> [kind](#module\_kind) | ../../modules/kind | n/a |
| <a name="module_loki"></a> [loki](#module\_loki) | ../../modules/loki | n/a |
| <a name="module_metrics"></a> [metrics](#module\_metrics) | ../../modules/metrics | n/a |
| <a name="module_mimir"></a> [mimir](#module\_mimir) | ../../modules/mimir | n/a |
| <a name="module_mysql"></a> [mysql](#module\_mysql) | ../../modules/mysql | n/a |
| <a name="module_postgresql"></a> [postgresql](#module\_postgresql) | ../../modules/postgresql | n/a |
| <a name="module_prometheus"></a> [prometheus](#module\_prometheus) | ../../modules/prometheus | n/a |
| <a name="module_redis"></a> [redis](#module\_redis) | ../../modules/redis | n/a |
| <a name="module_sql_mysql"></a> [sql\_mysql](#module\_sql\_mysql) | ../../modules/sql | n/a |
| <a name="module_sql_postgres"></a> [sql\_postgres](#module\_sql\_postgres) | ../../modules/sql | n/a |
| <a name="module_tempo"></a> [tempo](#module\_tempo) | ../../modules/tempo | n/a |

## Resources

| Name | Type |
|------|------|
| [random_password.sql_mysql_token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_password.sql_postgres_token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_pet.mysql_database](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.mysql_user](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.postgres_database](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.postgres_user](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.sql_mysql_host](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [random_pet.sql_postgres_host](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_alertmanager_url"></a> [alertmanager\_url](#output\_alertmanager\_url) | AlertManager URL |
| <a name="output_grafana_admin_password"></a> [grafana\_admin\_password](#output\_grafana\_admin\_password) | Grafana admin password |
| <a name="output_grafana_url"></a> [grafana\_url](#output\_grafana\_url) | Grafana URL |
| <a name="output_mysql_database"></a> [mysql\_database](#output\_mysql\_database) | MySQL database |
| <a name="output_mysql_password"></a> [mysql\_password](#output\_mysql\_password) | MySQL password |
| <a name="output_mysql_url"></a> [mysql\_url](#output\_mysql\_url) | MySQL URL |
| <a name="output_postgres_database"></a> [postgres\_database](#output\_postgres\_database) | Postgres database |
| <a name="output_postgres_password"></a> [postgres\_password](#output\_postgres\_password) | Postgres password |
| <a name="output_postgres_url"></a> [postgres\_url](#output\_postgres\_url) | Postgres URL |
| <a name="output_prometheus_url"></a> [prometheus\_url](#output\_prometheus\_url) | Prometheus URL |
| <a name="output_redis_password"></a> [redis\_password](#output\_redis\_password) | Redis password |
| <a name="output_redis_url"></a> [redis\_url](#output\_redis\_url) | Redis URL |
| <a name="output_sql_mysql_ingress_host"></a> [sql\_mysql\_ingress\_host](#output\_sql\_mysql\_ingress\_host) | SQL MySQL Ingress host |
| <a name="output_sql_mysql_ingress_url"></a> [sql\_mysql\_ingress\_url](#output\_sql\_mysql\_ingress\_url) | SQL MySQL Ingress URL |
| <a name="output_sql_mysql_token"></a> [sql\_mysql\_token](#output\_sql\_mysql\_token) | SQL MySQL token |
| <a name="output_sql_mysql_url"></a> [sql\_mysql\_url](#output\_sql\_mysql\_url) | SQL MySQL URL |
| <a name="output_sql_postgres_ingress_host"></a> [sql\_postgres\_ingress\_host](#output\_sql\_postgres\_ingress\_host) | SQL Postgres Ingress host |
| <a name="output_sql_postgres_ingress_url"></a> [sql\_postgres\_ingress\_url](#output\_sql\_postgres\_ingress\_url) | SQL Postgres Ingress URL |
| <a name="output_sql_postgres_token"></a> [sql\_postgres\_token](#output\_sql\_postgres\_token) | SQL Postgres token |
| <a name="output_sql_postgres_url"></a> [sql\_postgres\_url](#output\_sql\_postgres\_url) | SQL Postgres URL |
<!-- END_TF_DOCS -->
