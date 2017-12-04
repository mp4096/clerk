.DEFAULT_GOAL := build

build: fetch_dependencies ## Build
	go build -v github.com/mp4096/clerk/cmd/clerk

install: fetch_dependencies ## Build and install
	go install -v github.com/mp4096/clerk/cmd/clerk

xcompile_win: fetch_dependencies ## Cross-compile for Windows x64
	env GOOS=windows GOARCH=amd64 go build -v github.com/mp4096/clerk/cmd/clerk

fetch_dependencies: ## Fetch all dependencies
	go get -t ./...

fmt: ## Call go fmt in all directories
	go fmt .
	go fmt ./cmd/clerk/

delete_previews: ## Delete previews
	find . -type f -name 'clerk_preview_*' -delete

vet: ## Call go vet in all directories
	go vet .
	go vet ./cmd/clerk/

.PHONY: build install xcompile_win fmt delete_previews help vet fetch_dependencies

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
