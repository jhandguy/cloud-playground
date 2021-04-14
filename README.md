# DevOps Playground

[![CI](https://github.com/jhandguy/devops-playground/workflows/CI/badge.svg)](https://github.com/jhandguy/devops-playground/actions?query=workflow%3ACI)

A Playground to experiment with various DevOps tools and technologies.

## Tools

- Minikube
- LocalStack
- Prometheus
- Grafana
- AlertManager
- PushGateway
- Consul

## Technologies

- Terraform
- Kubernetes
- Helm

## Languages

- Golang
- YAML

## Architecture

```text
 -----------------------------------
|              [CONSUL]             |
|     -----------   -----------     |
|    | Dynamo DB | | S3 Bucket |    |
|     -----------   -----------     |
|          |             |          |
|         SDK           SDK         |
|          |             |          |
|      ----------   ----------      |   
|     |  dynamo  | |    s3    |     |
|      ----------   ----------      |
|            |         |            |
|           gRPC      gRPC          |
|            |         |            |
|            -----------            |
|           |  gateway  |           |
|            -----------            |
 -----------------------------------
         -------------------
        |  Ingress Gateway  |
         -------------------
                  |
                 REST
                  |
               -------
              |  cli  |
               -------
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