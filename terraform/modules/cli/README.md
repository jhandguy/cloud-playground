<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_helm"></a> [helm](#provider\_helm) | n/a |
| <a name="provider_null"></a> [null](#provider\_null) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [helm_release.cli](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [null_resource.logs](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_csi_enabled"></a> [csi\_enabled](#input\_csi\_enabled) | Enable CSI | `bool` | `false` | no |
| <a name="input_secrets"></a> [secrets](#input\_secrets) | Secrets | `map(string)` | `{}` | no |
| <a name="input_test_rounds"></a> [test\_rounds](#input\_test\_rounds) | Test rounds | `number` | `100` | no |
| <a name="input_vault_url"></a> [vault\_url](#input\_vault\_url) | Vault URL | `string` | `""` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
