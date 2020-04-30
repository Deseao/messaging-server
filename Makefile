COMPONENT=messaging-server
VERSION=$(shell cat VERSION)
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

MAIN=.
GO_FILES=$(shell find . -name \*.go)
MONITORED_FILES=$(GO_FILES) go.mod
UNIT_TEST_PACKAGES=. ./internal/... ./

$(COMPONENT): $(MONITORED_FILES)
	go build -o $(COMPONENT) $(MAIN)

.PHONY: ci-tools
ci-tools:
	GO111MODULE=off go get -u github.com/axw/gocov/...
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml

.PHONY: fmt
fmt: $(GO_FILES)
	go fmt ./...

.PHONY: test
test: $(GO_FILES)
	GIN_MODE=test go test -v -covermode=atomic -cover  $(UNIT_TEST_PACKAGES)

gocov.xml: $(GO_FILES)
	GIN_MODE=test gocov test -covermode=atomic $(UNIT_TEST_PACKAGES) | gocov-xml > gocov.xml

$(COVERAGE_FILE): $(GO_FILES)
	GIN_MODE=test go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) $(UNIT_TEST_PACKAGES)

$(COVERAGE_HTML): $(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

.PHONY: docker
docker:
	docker build -f Dockerfile -t local/hackathon/$(COMPONENT):$(VERSION) .
	docker tag local/hackathon/$(COMPONENT):$(VERSION) local/hackathon/$(COMPONENT):latest

.PHONY: clean
clean:
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML) $(COMPONENT) ./functional_test/logs/*.log
	find . -name \*.coverprofile -exec rm {} \;
	find . -name junit_\*_\*.xml -exec rm {} \;
	rm -f ./*.xml
.PHONY: lint
lint:
	golangci-lint run
