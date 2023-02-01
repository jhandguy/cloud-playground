FEATURE ?=

lint_helm:
	helm lint helm

lint:
	cargo fmt --all --check
	cargo clippy --features $(FEATURE) -- -D warnings
	cargo check --features $(FEATURE)

build:
	cargo build --features $(FEATURE) --bins --release --locked
	cargo test --features $(FEATURE)

test:
	cargo test --features $(FEATURE) -- --include-ignored

docker:
	docker buildx build --platform linux/arm64 --platform linux/amd64 --tag ghcr.io/jhandguy/cloud-playground/sql-$(FEATURE):latest --tag ghcr.io/jhandguy/cloud-playground/sql-$(FEATURE):1.0.0 --build-arg feature=$(FEATURE) --file Dockerfile . --push
