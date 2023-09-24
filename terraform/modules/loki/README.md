<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_helm"></a> [helm](#provider\_helm) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [helm_release.loki](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_alerting_rules"></a> [alerting\_rules](#input\_alerting\_rules) | Alerting rules | `list(string)` | `[]` | no |
| <a name="input_alertmanager_url"></a> [alertmanager\_url](#input\_alertmanager\_url) | AlertManager URL | `string` | n/a | yes |
| <a name="input_labels"></a> [labels](#input\_labels) | Labels | `list(string)` | <pre>[<br>  "level",<br>  "msg",<br>  "caller"<br>]</pre> | no |
| <a name="input_prometheus_enabled"></a> [prometheus\_enabled](#input\_prometheus\_enabled) | Enable Prometheus | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cluster_url"></a> [cluster\_url](#output\_cluster\_url) | Cluster URL |
<!-- END_TF_DOCS -->
