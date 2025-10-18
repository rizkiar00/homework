OAPI= oapi-codegen

.PHONY: deps gen build run tidy

deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go mod tidy

# generate server boilerplate into generated folder
gen:
	mkdir -p generated
	$(OAPI) -generate "types,chi-server" -o generated/api.gen.go api/openapi.yaml

# build the CLI/server
build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

tidy:
	go mod tidy
