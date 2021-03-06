version=0.1.0
export GOPROXY=direct
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

.PHONY: all dependencies docker clean

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  dependencies  - install dependencies"
	@echo "  build         - build the example"
	@echo "  clean         - clean the source directory"
	@echo "  lint          - lint the example"
	@echo "  fmt           - format the example"
	@echo "  docker        - build docker image for the example"

dependencies:
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/go-sif/sif@master
	@go get -d -v ./...

fmt:
	@go fmt ./...

clean:
	@rm -rf ./bin

lint:
	@echo "Running go vet"
	@go vet ./...
	@echo "Running golint"
	@go list ./... | xargs -L1 golint --set_exit_status

testenv:
	@echo "Downloading EDSM test files..."
	@mkdir -p testenv
	@cd testenv && curl -s https://www.edsm.net/dump/systemsWithCoordinates7days.json.gz | gunzip | tail -n +2 | head -n -1 | sed 's/,$$//' | sed 's/^....//' | split --additional-suffix .jsonl -l 50000
	@echo "Finished downloading EDSM test files."

build: lint
	@echo "Building sif docker swarm example..."
	@go build -a -ldflags="-w -s" -o bin/example.bin ./...
	@go mod tidy
	@echo "Finished building sif docker swarm example."

docker:
	@echo "Building sif docker swarm example image..."
	@docker build -t github.com/go-sif/sif-example-docker-swarm/example:latest .
	@echo "Finished building sif docker swarm example image."
