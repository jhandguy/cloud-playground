# Cloud Playground

[![CI](https://github.com/jhandguy/cloud-playground/actions/workflows/ci.yaml/badge.svg)](https://github.com/jhandguy/cloud-playground/actions/workflows/ci.yaml)

A Playground to experiment with various Cloud tools and technologies.

## Install Required Tools

```shell
brew install protobuf protoc-gen-go protoc-gen-go-grpc kind terraform k6
```

## Create Infrastructure

| Environment          | Command                                                         |
|----------------------|-----------------------------------------------------------------|
| Consul               | `make setup ENVIRONMENT=consul`                                 |
| Nginx                | `make setup ENVIRONMENT=nginx`                                  |
| Nginx (ArgoRollouts) | `make setup ENVIRONMENT=nginx TF_VAR_argorollouts_enabled=true` |
| Redis                | `make setup ENVIRONMENT=redis`                                  |

## Run Tests

| Environment             | Command                            |
|-------------------------|------------------------------------|
| Consul                  | `make go_test ENVIRONMENT=consul`  |
| Nginx                   | `make go_test ENVIRONMENT=nginx`   |
| Redis                   | `make rust_test ENVIRONMENT=redis` |

## Run Load Tests

| Environment             | Command                            |
|-------------------------|------------------------------------|
| Consul                  | `make go_load ENVIRONMENT=consul`  |
| Nginx                   | `make go_load ENVIRONMENT=nginx`   |
| Redis                   | `make rust_load ENVIRONMENT=redis` |

## Destroy Infrastructure

| Environment             | Command                            |
|-------------------------|------------------------------------|
| Consul                  | `make teardown ENVIRONMENT=consul` |
| Nginx                   | `make teardown ENVIRONMENT=nginx`  |
| Redis                   | `make teardown ENVIRONMENT=redis`  |
