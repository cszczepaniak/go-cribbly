.PHONY: serve
serve:
	go run cmd/local/main.go

.PHONY: build
build:
	./scripts/build-lamda.sh