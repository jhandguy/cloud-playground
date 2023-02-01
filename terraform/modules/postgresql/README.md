<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_helm"></a> [helm](#provider\_helm) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [helm_release.postgresql](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [random_password.postgres_password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_password.user_password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_database_name"></a> [database\_name](#input\_database\_name) | Database name | `string` | n/a | yes |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_port"></a> [node\_port](#input\_node\_port) | Node port | `number` | n/a | yes |
| <a name="input_user_name"></a> [user\_name](#input\_user\_name) | User name | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cluster_url"></a> [cluster\_url](#output\_cluster\_url) | Cluster URL |
| <a name="output_database_name"></a> [database\_name](#output\_database\_name) | Database name |
| <a name="output_postgres_password"></a> [postgres\_password](#output\_postgres\_password) | Postgres password |
| <a name="output_url"></a> [url](#output\_url) | URL |
| <a name="output_user_name"></a> [user\_name](#output\_user\_name) | User name |
| <a name="output_user_password"></a> [user\_password](#output\_user\_password) | User password |
<!-- END_TF_DOCS -->
