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

## Automation

```shell
make all
```

### Create Infrastructure

```shell
make setup
```

### Run Tests

```shell
make test
```

### Destroy Infrastructure

```shell
make teardown
```