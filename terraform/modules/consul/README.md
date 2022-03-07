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
| [helm_release.consul](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_ingress_gateway_port"></a> [ingress\_gateway\_port](#input\_ingress\_gateway\_port) | Ingress Gateway port | `number` | `8080` | no |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_port"></a> [node\_port](#input\_node\_port) | Node port | `number` | n/a | yes |
| <a name="input_node_ports"></a> [node\_ports](#input\_node\_ports) | Node ports | `map(number)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_ingress_gateway_cluster_url"></a> [ingress\_gateway\_cluster\_url](#output\_ingress\_gateway\_cluster\_url) | Ingress Gateway cluster URL |
| <a name="output_ingress_gateway_port"></a> [ingress\_gateway\_port](#output\_ingress\_gateway\_port) | Ingress Gateway port |
| <a name="output_ingress_gateway_urls"></a> [ingress\_gateway\_urls](#output\_ingress\_gateway\_urls) | Ingress Gateway URLs |
| <a name="output_url"></a> [url](#output\_url) | URL |
<!-- END_TF_DOCS -->
