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
| [helm_release.prometheus](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [random_password.admin_password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_alertmanager_node_port"></a> [alertmanager\_node\_port](#input\_alertmanager\_node\_port) | AlertManager node port | `number` | n/a | yes |
| <a name="input_grafana_dashboards"></a> [grafana\_dashboards](#input\_grafana\_dashboards) | Grafana dashboards | `list(string)` | `[]` | no |
| <a name="input_grafana_datasources"></a> [grafana\_datasources](#input\_grafana\_datasources) | Grafana datasources | `list(string)` | `[]` | no |
| <a name="input_grafana_node_port"></a> [grafana\_node\_port](#input\_grafana\_node\_port) | Grafana node port | `number` | n/a | yes |
| <a name="input_mimir_url"></a> [mimir\_url](#input\_mimir\_url) | Mimir URL | `string` | `""` | no |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_prometheus_node_port"></a> [prometheus\_node\_port](#input\_prometheus\_node\_port) | Prometheus node port | `number` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_alertmanager_cluster_url"></a> [alertmanager\_cluster\_url](#output\_alertmanager\_cluster\_url) | AlertManager Cluster URL |
| <a name="output_alertmanager_url"></a> [alertmanager\_url](#output\_alertmanager\_url) | AlertManager URL |
| <a name="output_grafana_admin_password"></a> [grafana\_admin\_password](#output\_grafana\_admin\_password) | Grafana admin password |
| <a name="output_grafana_url"></a> [grafana\_url](#output\_grafana\_url) | Grafana URL |
| <a name="output_prometheus_cluster_url"></a> [prometheus\_cluster\_url](#output\_prometheus\_cluster\_url) | Prometheus URL |
| <a name="output_prometheus_url"></a> [prometheus\_url](#output\_prometheus\_url) | Prometheus URL |
<!-- END_TF_DOCS -->
