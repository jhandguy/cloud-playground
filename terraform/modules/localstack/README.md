<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_helm"></a> [helm](#provider\_helm) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_dynamodb_table.dynamo](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |
| [aws_s3_bucket.s3](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [helm_release.localstack](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [random_id.buckets](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id) | resource |
| [random_id.tables](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_dynamo_tables"></a> [aws\_dynamo\_tables](#input\_aws\_dynamo\_tables) | AWS DynamoDB tables | `list(string)` | `[]` | no |
| <a name="input_aws_s3_buckets"></a> [aws\_s3\_buckets](#input\_aws\_s3\_buckets) | AWS S3 buckets | `list(string)` | `[]` | no |
| <a name="input_node_ip"></a> [node\_ip](#input\_node\_ip) | Node ip | `string` | n/a | yes |
| <a name="input_node_port"></a> [node\_port](#input\_node\_port) | Node port | `number` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_aws_dynamo_endpoint"></a> [aws\_dynamo\_endpoint](#output\_aws\_dynamo\_endpoint) | AWS DynamoDB endpoint |
| <a name="output_aws_dynamo_tables"></a> [aws\_dynamo\_tables](#output\_aws\_dynamo\_tables) | AWS DynamoDB tables |
| <a name="output_aws_s3_buckets"></a> [aws\_s3\_buckets](#output\_aws\_s3\_buckets) | AWS S3 buckets |
| <a name="output_aws_s3_endpoint"></a> [aws\_s3\_endpoint](#output\_aws\_s3\_endpoint) | AWS S3 endpoint |
<!-- END_TF_DOCS -->
