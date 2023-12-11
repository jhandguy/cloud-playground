<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_null"></a> [null](#provider\_null) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [null_resource.cluster](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [random_shuffle.node_ports](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/shuffle) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_name"></a> [cluster\_name](#input\_cluster\_name) | Cluster name | `string` | `"kind"` | no |
| <a name="input_node_image"></a> [node\_image](#input\_node\_image) | Node image | `string` | `"v1.28.0"` | no |
| <a name="input_node_ports"></a> [node\_ports](#input\_node\_ports) | Node ports | `list(string)` | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cluster_context"></a> [cluster\_context](#output\_cluster\_context) | Cluster context |
| <a name="output_node_ip"></a> [node\_ip](#output\_node\_ip) | Node ip |
| <a name="output_node_ports"></a> [node\_ports](#output\_node\_ports) | Node ports |
<!-- END_TF_DOCS -->