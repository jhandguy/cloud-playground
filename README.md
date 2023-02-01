# Cloud Playground

[![CI](https://github.com/jhandguy/cloud-playground/actions/workflows/ci.yaml/badge.svg)](https://github.com/jhandguy/cloud-playground/actions/workflows/ci.yaml)

A Playground to experiment with various Cloud tools and technologies.

## Tools

- Kind
- LocalStack
- Prometheus
- Grafana
- Promtail
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
- ArgoRollouts
- MetricsServer
- PostgreSQL
- MySQL
- Redis

## Technologies

- Terraform
- Kubernetes
- Helm

## Languages

- Golang
- Rust
- YAML
- HCL
- JavaScript

## Protocols

- gRPC
- REST

## Install Required Tools

```shell
brew install protobuf protoc-gen-go protoc-gen-go-grpc kind terraform k6
```

## Create Infrastructure

### Consul

```shell
make setup ENVIRONMENT=consul
```

### Nginx

```shell
make setup ENVIRONMENT=nginx
```

### Nginx + ArgoRollouts

```shell
make setup ENVIRONMENT=nginx TF_VAR_argorollouts_enabled=true
```

### Redis

```shell
make setup ENVIRONMENT=redis
```

## Run Tests

### Consul

```shell
make go_test ENVIRONMENT=consul
```

### Nginx

```shell
make go_test ENVIRONMENT=nginx
```

## Run Load Tests

### Consul

```shell
make go_load ENVIRONMENT=consul
```

### Nginx

```shell
make go_load ENVIRONMENT=nginx
```

## Destroy Infrastructure

### Consul

```shell
make teardown ENVIRONMENT=consul
```

### Nginx

```shell
make teardown ENVIRONMENT=nginx
```

### Redis

```shell
make teardown ENVIRONMENT=redis
```
