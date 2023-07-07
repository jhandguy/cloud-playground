ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
FEATURE ?=

lint_helm:
	helm lint helm

lint:
	cargo fmt --all --check
	cargo clippy --features $(FEATURE) -- -D warnings
	cargo check --features $(FEATURE)

build:
	cargo build --features $(FEATURE) --bins --release --locked --target x86_64-unknown-linux-musl
	cargo test --features $(FEATURE)

test:
	cargo test --features $(FEATURE) -- --include-ignored

load: export SQL_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json sql_$(FEATURE)_token | jq -r .)
load: export SQL_HOST=$(shell terraform -chdir=$(CHDIR) output -json sql_$(FEATURE)_ingress_host | jq -r .)
load: export SQL_URL=$(shell terraform -chdir=$(CHDIR) output -json sql_$(FEATURE)_ingress_url | jq -r .)
load:
	k6 run k6/script.js

docker:
	docker buildx build --platform linux/arm64 --platform linux/amd64 --tag ghcr.io/jhandguy/cloud-playground/sql-$(FEATURE):latest --tag ghcr.io/jhandguy/cloud-playground/sql-$(FEATURE):1.0.0 --build-arg feature=$(FEATURE) --file Dockerfile . --push
