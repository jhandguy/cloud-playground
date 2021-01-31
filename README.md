# DevOps Playground

[![CI](https://github.com/jhandguy/devops-playground/workflows/CI/badge.svg)](https://github.com/jhandguy/devops-playground/actions?query=workflow%3ACI)

A Playground to experiment with various DevOps tools and technologies.

## Tools

- Minikube
- LocalStack

## Technologies

- Terraform
- Kubernetes
- Helm

## Languages

- Golang
- YAML

## Architecture

```text
         -------
        |  cli  |
         -------
            |
           REST
            |
       -----------
      |  gateway  |
       -----------
       |         |
      gRPC      gRPC
       |         |
 ----------   ----------   
|  dynamo  | |    s3    |
 ----------   ----------
```

## Create Infrastructure

```shell
chmod +x create && ./create
```

## Run Tests

```shell
chmod +x test && ./test
```

## Destroy Infrastructure

```shell
chmod +x destroy && ./destroy
```