<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.1.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cli"></a> [cli](#module\_cli) | ./modules/cli | n/a |
| <a name="module_consul"></a> [consul](#module\_consul) | ./modules/consul | n/a |
| <a name="module_csi"></a> [csi](#module\_csi) | ./modules/csi | n/a |
| <a name="module_dynamo"></a> [dynamo](#module\_dynamo) | ./modules/dynamo | n/a |
| <a name="module_gateway"></a> [gateway](#module\_gateway) | ./modules/gateway | n/a |
| <a name="module_localstack"></a> [localstack](#module\_localstack) | ./modules/localstack | n/a |
| <a name="module_loki"></a> [loki](#module\_loki) | ./modules/loki | n/a |
| <a name="module_metrics"></a> [metrics](#module\_metrics) | ./modules/metrics | n/a |
| <a name="module_minikube"></a> [minikube](#module\_minikube) | ./modules/minikube | n/a |
| <a name="module_prometheus"></a> [prometheus](#module\_prometheus) | ./modules/prometheus | n/a |
| <a name="module_pushgateway"></a> [pushgateway](#module\_pushgateway) | ./modules/pushgateway | n/a |
| <a name="module_s3"></a> [s3](#module\_s3) | ./modules/s3 | n/a |
| <a name="module_vault"></a> [vault](#module\_vault) | ./modules/vault | n/a |

## Resources

| Name | Type |
|------|------|
| [random_password.dynamo_token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_password.gateway_token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_password.s3_token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_access_key_id"></a> [aws\_access\_key\_id](#input\_aws\_access\_key\_id) | AWS access key id | `string` | n/a | yes |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region | `string` | n/a | yes |
| <a name="input_aws_secret_access_key"></a> [aws\_secret\_access\_key](#input\_aws\_secret\_access\_key) | AWS secret access key | `string` | n/a | yes |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_alertmanager_url"></a> [alertmanager\_url](#output\_alertmanager\_url) | AlertManager URL |
| <a name="output_aws_dynamo_endpoint"></a> [aws\_dynamo\_endpoint](#output\_aws\_dynamo\_endpoint) | Dynamo endpoint |
| <a name="output_aws_dynamo_table"></a> [aws\_dynamo\_table](#output\_aws\_dynamo\_table) | Dynamo table |
| <a name="output_aws_s3_bucket"></a> [aws\_s3\_bucket](#output\_aws\_s3\_bucket) | S3 bucket |
| <a name="output_aws_s3_endpoint"></a> [aws\_s3\_endpoint](#output\_aws\_s3\_endpoint) | S3 endpoint |
| <a name="output_canary_gateway_url"></a> [canary\_gateway\_url](#output\_canary\_gateway\_url) | Canary Gateway URL |
| <a name="output_consul_url"></a> [consul\_url](#output\_consul\_url) | Consul URL |
| <a name="output_dynamo_token"></a> [dynamo\_token](#output\_dynamo\_token) | Dynamo token |
| <a name="output_dynamo_url"></a> [dynamo\_url](#output\_dynamo\_url) | Dynamo URL |
| <a name="output_gateway_token"></a> [gateway\_token](#output\_gateway\_token) | Gateway token |
| <a name="output_grafana_admin_password"></a> [grafana\_admin\_password](#output\_grafana\_admin\_password) | Grafana admin password |
| <a name="output_grafana_url"></a> [grafana\_url](#output\_grafana\_url) | Grafana URL |
| <a name="output_ingress_gateway_url"></a> [ingress\_gateway\_url](#output\_ingress\_gateway\_url) | Gateway Ingress URL |
| <a name="output_prod_gateway_url"></a> [prod\_gateway\_url](#output\_prod\_gateway\_url) | Prod Gateway URL |
| <a name="output_prometheus_url"></a> [prometheus\_url](#output\_prometheus\_url) | Prometheus URL |
| <a name="output_pushgateway_url"></a> [pushgateway\_url](#output\_pushgateway\_url) | PushGateway URL |
| <a name="output_s3_token"></a> [s3\_token](#output\_s3\_token) | S3 token |
| <a name="output_s3_url"></a> [s3\_url](#output\_s3\_url) | S3 URL |
| <a name="output_vault_root_token"></a> [vault\_root\_token](#output\_vault\_root\_token) | Vault root token |
| <a name="output_vault_url"></a> [vault\_url](#output\_vault\_url) | Vault URL |
<!-- END_TF_DOCS -->
