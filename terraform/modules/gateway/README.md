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
| [helm_release.gateway](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_argorollouts_enabled"></a> [argorollouts\_enabled](#input\_argorollouts\_enabled) | Enable ArgoRollouts | `bool` | `false` | no |
| <a name="input_consul_enabled"></a> [consul\_enabled](#input\_consul\_enabled) | Enable Consul | `bool` | `false` | no |
| <a name="input_csi_enabled"></a> [csi\_enabled](#input\_csi\_enabled) | Enable CSI | `bool` | `false` | no |
| <a name="input_ingress_gateway_port"></a> [ingress\_gateway\_port](#input\_ingress\_gateway\_port) | Ingress Gateway port | `number` | `8080` | no |
| <a name="input_ingress_host"></a> [ingress\_host](#input\_ingress\_host) | Ingress host | `string` | n/a | yes |
| <a name="input_max_replicas"></a> [max\_replicas](#input\_max\_replicas) | Maximum replicas | `number` | `1` | no |
| <a name="input_min_replicas"></a> [min\_replicas](#input\_min\_replicas) | Minimum replicas | `number` | `1` | no |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_ports"></a> [node\_ports](#input\_node\_ports) | Node ports | `map(number)` | `{}` | no |
| <a name="input_prometheus_enabled"></a> [prometheus\_enabled](#input\_prometheus\_enabled) | Enable Prometheus | `bool` | `false` | no |
| <a name="input_prometheus_url"></a> [prometheus\_url](#input\_prometheus\_url) | Prometheus URL | `string` | `""` | no |
| <a name="input_secrets"></a> [secrets](#input\_secrets) | Secrets | `map(string)` | `{}` | no |
| <a name="input_vault_url"></a> [vault\_url](#input\_vault\_url) | Vault URL | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_host"></a> [host](#output\_host) | Host |
| <a name="output_urls"></a> [urls](#output\_urls) | URLs |
<!-- END_TF_DOCS -->
