all: setup test teardown

setup: setup_minikube setup_terraform

setup_minikube:
	minikube start $(shell if [[ $OSTYPE != "linux-gnu"* ]]; then echo "--vm=true"; fi)

setup_terraform:
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -var="node_ip=$(shell minikube ip)" -out=tfplan
	terraform -chdir=terraform apply tfplan
	rm terraform/tfplan

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