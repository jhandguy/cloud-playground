export AWS_REGION=eu-central-1
export AWS_ACCESS_KEY_ID=aws-access-key-id
export AWS_SECRET_ACCESS_KEY=aws-secret-access-key
export AWS_DEFAULT_REGION=$(AWS_REGION)
export TF_VAR_aws_region=$(AWS_REGION)
export TF_VAR_aws_access_key_id=$(AWS_ACCESS_KEY_ID)
export TF_VAR_aws_secret_access_key=$(AWS_SECRET_ACCESS_KEY)

ENVIRONMENT ?=
CHDIR = terraform/environments/$(ENVIRONMENT)

go_ci: lint_terraform lint_helm setup go_compile go_build go_test go_load teardown

rust_ci: lint_terraform lint_helm setup rust_build rust_test teardown

setup:
	terraform -chdir=$(CHDIR) init
	terraform -chdir=$(CHDIR) plan -out=tfplan
	terraform -chdir=$(CHDIR) apply tfplan
	rm $(CHDIR)/tfplan

teardown:
	terraform -chdir=$(CHDIR) destroy -auto-approve

go_compile:
	make -j compile_s3 compile_dynamo compile_gateway

compile_%:
	make -C $* compile

format:
	terraform fmt -recursive
	terraform-docs markdown table $(CHDIR) --output-file README.md --recursive --recursive-path ../../modules

lint_terraform:
	terraform fmt -recursive -check

lint_helm:
	make -j lint_helm_s3 lint_helm_dynamo lint_helm_gateway lint_helm_cli lint_helm_sql

lint_helm_%:
	make -C $* lint_helm

go_lint: lint_s3 lint_dynamo lint_gateway lint_cli

rust_lint:
	lint_sql FEATURE=postgres
	lint_sql FEATURE=mysql

lint_%:
	make -C $* lint

go_build:
	make -j build_s3 build_dynamo build_gateway build_cli

rust_build:
	make build_sql FEATURE=postgres
	make build_sql FEATURE=mysql

build_%:
	make -C $* build

go_test: test_s3 test_dynamo test_gateway test_cli

rust_test:
	make test_sql FEATURE=postgres
	make test_sql FEATURE=mysql

test_%:
	make -C $* test HTTP_PORT=8080 GRPC_PORT=8080 METRICS_PORT=9090

go_load: load_s3 load_dynamo load_gateway

load_%:
	make -C $* load

update:
	terraform -chdir=$(CHDIR) init -upgrade
	terraform -chdir=$(CHDIR) providers lock
	make -j update_s3 update_dynamo update_gateway update_cli

update_%:
	make -C $* update

go_docker: docker_s3 docker_dynamo docker_gateway docker_cli

rust_docker:
	make docker_sql FEATURE=postgres
	make docker_sql FEATURE=mysql

docker_%:
	make -C $* docker
