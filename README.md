# DevOps Playground

[![CI](https://github.com/jhandguy/devops-playground/workflows/CI/badge.svg)](https://github.com/jhandguy/devops-playground/actions?query=workflow%3ACI)

A Playground to experiment with various DevOps tools and technologies.

## Tools

- Minikube
- LocalStack
- Prometheus
- Grafana
- Loki
- Tempo
- AlertManager
- PushGateway
- Consul
- Vault
- CSI
- K6
- IngressNGINX
- CertManager

## Technologies

- Terraform
- Kubernetes
- Helm

## Languages

- Golang
- YAML
- HCL
- JavaScript

## Protocols

- gRPC
- REST

## Architecture

```text
 -----------------------------------
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
|         -----------------         |
|        |     gateway     |        |
|        | _______ _______ |        |
|        |  prod  | canary |        |
|         -----------------         |
|            ||       ||            |
|           50%       50%           |
|            ||       ||            |
 -----------------------------------
             -----------
            |  Ingress  |
             -----------
                  |
                 REST
                  |
               -------
              |  cli  |
               -------
```

### Install Required Tools

```shell
brew install protobuf protoc-gen-go protoc-gen-go-grpc minikube terraform k6
```

### Create Infrastructure

#### Consul

```shell
make setup ENVIRONMENT=consul
```

#### Nginx

```shell
make setup ENVIRONMENT=nginx
```

### Run Tests

#### Consul

```shell
make test ENVIRONMENT=consul
```

#### Nginx

```shell
make test ENVIRONMENT=nginx
```

### Run Load Tests

#### Consul

```shell
make load ENVIRONMENT=consul
```

#### Nginx

```shell
make load ENVIRONMENT=nginx
```

### Destroy Infrastructure

#### Consul

```shell
make teardown ENVIRONMENT=consul
```

#### Nginx

```shell
make teardown ENVIRONMENT=nginx
```
