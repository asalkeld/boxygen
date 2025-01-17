install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

test: generate-proto
	@echo Running Tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./pkg/...

generate-proto: install-tools
	@echo Generating Proto Sources
	@mkdir -p ./pkg/proto
	@protoc --go_out=./pkg/proto --go-grpc_out=./pkg/proto -I ./proto ./proto/*/**/*.proto

generate-mocks: install-tools generate-proto
	@echo Generating mocks
	@mkdir -p ./mocks/proto
	@go run github.com/golang/mock/mockgen github.com/nitrictech/boxygen/pkg/proto/builder/v1 Builder_AddServer,Builder_ConfigServer,Builder_CopyServer,Builder_RunServer > mocks/proto/mock.go

build: generate-proto
	@CGO_ENABLED=0 go build -tags containers_image_openpgp -o bin/boxygen -ldflags="-extldflags=-static" ./cmd/docker/run.go

build-docker:
	@DOCKER_BUILDKIT=1 docker build . -f docker/docker/Dockerfile -t nitrictech/boxygen-dockerfile

sourcefiles := $(shell find . -type f -name "*.go")

license-header-check:
	@echo Checking License Headers to Source Files
	@addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)

license-check: install-tools license-check-dev license-check-aws license-check-gcp license-check-azure
	@echo Checking OSS Licenses