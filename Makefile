ci: compile lint build setup test teardown

compile:
	make -j compile_s3 compile_dynamo compile_gateway

compile_s3:
	make -C s3 compile

compile_dynamo:
	make -C dynamo compile

compile_gateway:
	make -C gateway compile

lint:
	make -j lint_terraform lint_helm lint_golang

lint_terraform:
	terraform fmt -recursive -check

lint_helm:
	helm lint s3/helm dynamo/helm

lint_golang:
	make -C s3 lint
	make -C dynamo lint
	make -C gateway lint
	make -C cli lint

build:
	make -j build_s3 build_dynamo build_gateway build_cli

build_s3:
	make -C s3 build

build_dynamo:
	make -C dynamo build

build_gateway:
	make -C gateway build

build_cli:
	make -C cli build

setup: setup_minikube setup_terraform

setup_minikube:
	minikube start $(shell if [ $$(uname) != "Linux" ]; then echo "--vm=true"; fi)

setup_terraform:
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -var="node_ip=$(shell minikube ip)" -out=tfplan
	terraform -chdir=terraform apply tfplan
	rm terraform/tfplan

test:
	make -j test_s3 test_dynamo test_gateway
	make test_cli ROUNDS=10

test_s3:
	make -C s3 test GRPC_PORT=8080 METRICS_PORT=9090

test_dynamo:
	make -C dynamo test GRPC_PORT=8081 METRICS_PORT=9091

test_gateway:
	make -C gateway test HTTP_PORT=8082 METRICS_PORT=9092

test_cli:
	make -C cli test

teardown: teardown_terraform teardown_minikube

teardown_terraform:
	terraform -chdir=terraform destroy -var="node_ip=$(shell minikube ip)" -auto-approve

teardown_minikube:
	minikube stop
	minikube delete

update:
	make -j update_s3 update_dynamo update_gateway update_cli

update_s3:
	make -C s3 update

update_dynamo:
	make -C dynamo update

update_gateway:
	make -C gateway update

update_cli:
	make -C cli update

docker:
	make -j docker_s3 docker_dynamo docker_gateway docker_cli

docker_s3:
	make -C s3 docker

docker_dynamo:
	make -C dynamo docker

docker_gateway:
	make -C gateway docker

docker_cli:
	make -C cli docker