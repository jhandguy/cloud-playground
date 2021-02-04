all: compile lint build setup test teardown

setup: setup_minikube setup_terraform

setup_minikube:
	minikube start $(shell if [[ $OSTYPE != "linux-gnu"* ]]; then echo "--vm=true"; fi)

setup_terraform:
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -var="node_ip=$(shell minikube ip)" -out=tfplan
	terraform -chdir=terraform apply tfplan
	rm terraform/tfplan

compile:
	make -j compile_s3 compile_dynamo

compile_s3:
	make -C s3 compile

compile_dynamo:
	make -C dynamo compile

lint:
	make -j lint_terraform lint_helm lint_golang

lint_terraform:
	terraform fmt -recursive -check

lint_helm:
	helm lint s3/helm dynamo/helm

lint_golang:
	make -C s3 lint
	make -C dynamo lint

build:
	make -j build_s3 build_dynamo

build_s3:
	make -C s3 build

build_dynamo:
	make -C dynamo build

test:
	make -j test_s3 test_dynamo

test_s3:
	make -C s3 test

test_dynamo:
	make -C dynamo test

teardown: teardown_terraform teardown_minikube

teardown_terraform:
	terraform -chdir=terraform destroy -var="node_ip=$(shell minikube ip)" -auto-approve

teardown_minikube:
	minikube stop
	minikube delete