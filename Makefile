VERSION ?= 0.0.1
# Image URL to use all building/pushing image targets
IMG ?= cuttingedge1109/jsonschema-validation-webhook:v$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: lint

# Test
test:
	go test ./...

# Lint
lint:
	golangci-lint run --timeout=10m ./...


run: 
	cd ../
	revel run -a github.com/cuttingedge1109/jsonschema-validation-webhook

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build: test
	docker build -t ${IMG} .

# Push the docker image
docker-push:
	docker push ${IMG}

# Download controller-gen locally if necessary
REVEL = $(shell pwd)/bin/controller-gen
revel:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	REVEL_TMP_DIR=$$(mktemp -d) ;\
	cd $$REVEL_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/revel/cmd@v1.0.3 ;\
	rm -rf $$REVEL_TMP_DIR ;\
	}
REVEL=$(GOBIN)/revel
else
REVEL=$(shell which revel)
endif

bump-chart:
	sed -i "s/^version:.*/version:  $(VERSION)/" charts/validation-webhook/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion:  $(VERSION)/" charts/validation-webhook/Chart.yaml
	sed -i "s/tag:.*/tag:  v$(VERSION)/" charts/validation-webhook/values.yaml

helm-lint:
	if [ -d charts ]; then helm lint charts/validation-webhook; fi
