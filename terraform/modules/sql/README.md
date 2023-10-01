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
| [helm_release.sql](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_feature"></a> [feature](#input\_feature) | Feature | `string` | n/a | yes |
| <a name="input_ingress_host"></a> [ingress\_host](#input\_ingress\_host) | Ingress host | `string` | n/a | yes |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_ports"></a> [node\_ports](#input\_node\_ports) | Node ports | `tuple([number, number])` | n/a | yes |
| <a name="input_prometheus_enabled"></a> [prometheus\_enabled](#input\_prometheus\_enabled) | Enable Prometheus | `bool` | `true` | no |
| <a name="input_rate_limit_requests"></a> [rate\_limit\_requests](#input\_rate\_limit\_requests) | Rate limit requests | `number` | `0` | no |
| <a name="input_replicas"></a> [replicas](#input\_replicas) | Replicas | `number` | `1` | no |
| <a name="input_secrets"></a> [secrets](#input\_secrets) | Secrets | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_url"></a> [url](#output\_url) | URL |
<!-- END_TF_DOCS -->
