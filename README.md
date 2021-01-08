# DevOps Playground

[![CI](https://github.com/jhandguy/devops-playground/workflows/CI/badge.svg)](https://github.com/jhandguy/devops-playground/actions?query=workflow%3ACI)

A Playground to experiment with various DevOps tools and technologies.

## Tools

- Minikube
- LocalStack
- Jenkins

## Technologies

- Terraform
- Kubernetes
- Helm

## Create Infrastructure

### MacOS

```shell
chmod +x create && ./create --vm=true
```

### Linux

```shell
chmod +x create && ./create
```

## Run Tests

```shell
export AWS_S3_ENDPOINT=$(terraform -chdir=terraform output -json localstack | jq -r .)
export AWS_S3_BUCKET=$(terraform -chdir=terraform output -json bucket | jq -r .)
export S3_URL=$(terraform -chdir=terraform output -json s3 | jq -r .)
(cd s3 && go test ./... -cover -race)
```

## Destroy Infrastructure

```shell
chmod +x destroy && ./destroy
```