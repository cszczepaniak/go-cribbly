.PHONY: serve
serve:
	go run cmd/local/main.go

.PHONY: build
build:
	./scripts/build-lamda.sh

.PHONY: generate-storage
generate-storage:
	go run tool/storage/generate.go