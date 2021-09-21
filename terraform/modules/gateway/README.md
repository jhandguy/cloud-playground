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
| <a name="input_ingress_gateway_port"></a> [ingress\_gateway\_port](#input\_ingress\_gateway\_port) | Ingress Gateway port | `number` | n/a | yes |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_ports"></a> [node\_ports](#input\_node\_ports) | Node ports | `map(number)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_urls"></a> [urls](#output\_urls) | URLs |
<!-- END_TF_DOCS -->