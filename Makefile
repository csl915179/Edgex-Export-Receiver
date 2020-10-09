.PHONY: build clean test run docker

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/edgex-export-receiver
.PHONY: $(MICROSERVICES)


VERSION= master

GOFLAGS=-ldflags "-X github.com/edgexfoundry/edgex-ui-go.Version=$(VERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/edgex-export-receiver:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

clean:
	rm -f $(MICROSERVICES)

test:
	GO111MODULE=on go test -coverprofile=coverage.out ./...
	GO111MODULE=on go vet ./...

prepare:

run:
	cd cmd && ./edgex-export-receiver
docker: $(DOCKERS)

#docker_edgex_ui_go:
#	docker build --label "git_sha=$(GIT_SHA)" -t edgexfoundry/docker-edgex-ui-go:$(VERSION) .
