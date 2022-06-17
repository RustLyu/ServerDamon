PACKAGE  							= ${shell pwd | rev | cut -f1 -d'/' - | rev}
DATE    							?= $(shell date +%Y-%m-%d_%I:%M:%S%p)
GITHASH 							= $(shell git rev-parse HEAD)
VERSION 							?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
										cat $(PACKAGE)/.version 2> /dev/null || echo v0)
PKGS     							= $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS 							= $(shell env GO111MODULE=on $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN      							= $(GOPATH)/bin
HOST_PORT							= 5001
DOCKER_EXPOSED_PORT						= 5001
DOCKER_BUILD_CONTEXT						= .
DOCKER_FILE_PATH						= build/Dockerfile
GO      							= go
TIMEOUT 							= 300
V 								= 0
Q 								= $(if $(filter 1,$V),,@)
M 								= $(shell printf "\033[34;1m▶\033[0m")
os 								= $(shell uname)

# Api Specifications
API_DIR=api/specification

# Generated Code
GEN=gen
SERVICE_OUT=${GEN}/service
GRPC_OUT=${SERVICE_OUT}/grpc/pb

#this sets the right cert directory WSL(bottom) vs MacOS (top)
ifeq ($(os),$(filter $(os),Darwin Linux))
DOCKER_MNT_PATH				= ~/opt/certificates:/opt/certificates
else
DOCKER_MNT_PATH				= "~\\opt\\certificates:/opt/certificates"
endif
export GOPRIVATE="gitlab.com/mwdavisii*" #bypass validaton on private repos
export GO111MODULE=on
export MY_POD_NAMESPACE=development

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(REPOSITORY)…)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(REPOSITORY) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: REPOSITORY=golang.org/x/lint/golint

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: REPOSITORY=github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: REPOSITORY=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: REPOSITORY=github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: REPOSITORY=github.com/tebeka/go2xunit


########################################################################################################################
##########                                                                                                    ##########
########## (~˘▾˘)~  (~˘▾˘)~  (~˘▾˘)~  (~˘▾˘)~  (~˘▾˘)~  RECIPES  ~(˘▾˘~)  ~(˘▾˘~)  ~(˘▾˘~)  ~(˘▾˘~)  ~(˘▾˘~)  ##########
##########                                                                                                    ##########
########################################################################################################################

######################################################
#Build tools (for the CI/CD Pipelines)
######################################################

.DEFAULT_GOAL := build
build: all

.PHONY: all
all: fmt test $(BIN) 
	$(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build -tags release -ldflags '-X main.GitComHash=$(GITHASH) -X main.BuildStamp=$(DATE)' -o bin/application cmd/server/main.go

.PHONY: build-only
build-only: 
	$(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build -tags release -ldflags '-X main.GitComHash=$(GITHASH) -X main.BuildStamp=$(DATE)' -o bin/application cmd/server/main.go

build-only-musl: 
	$(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build -tags musl -ldflags '-X main.GitComHash=$(GITHASH) -X main.BuildStamp=$(DATE)' -o bin/application cmd/server/main.go


docker-pipeline:
	$(info $(M) building container...) @ ## Build docker container
	docker build $(DOCKER_BUILD_ARGS) -t $(PACKAGE):$(GITHASH) $(DOCKER_BUILD_CONTEXT) -f $(DOCKER_FILE_PATH)
	@DOCKER_MAJOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f1) ; \
	DOCKER_MINOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f2) ; \
	if [ $$DOCKER_MAJOR -eq 1 ] && [ $$DOCKER_MINOR -lt 10 ] ; then \
		echo docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
		docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
	else \
		echo docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ;\
		docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ; \
	fi

docker: all
	$(info $(M) building container...) @ ## Build docker container
	docker build $(DOCKER_BUILD_ARGS) -t $(PACKAGE):$(GITHASH) $(DOCKER_BUILD_CONTEXT) -f $(DOCKER_FILE_PATH)
	@DOCKER_MAJOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f1) ; \
	DOCKER_MINOR=$(shell docker -v | sed -e 's/.*version //' -e 's/,.*//' | cut -d\. -f2) ; \
	if [ $$DOCKER_MAJOR -eq 1 ] && [ $$DOCKER_MINOR -lt 10 ] ; then \
		echo docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
		docker tag -f $(PACKAGE):$(GITHASH) $(IMAGE):latest ;\
	else \
		echo docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ;\
		docker tag $(PACKAGE):$(GITHASH) $(PACKAGE):latest ; \
	fi

.PHONY: version
version:
	@echo $(VERSION)

######################################################
#Environment Tools (deps, certs, etc. for development)
######################################################
ENV_TARGETS := local-env-up local-env-down certs
.PHONY: $(ENV_TARGETS)
local-env-up: 
	cd test/env/local; \
	docker-compose up --detach --force-recreate --remove-orphans
local-env-down: 
	cd test/env/local; \
	docker-compose down --remove-orphans

certs:
	mkdir -p ~/opt/certificates/
	openssl genrsa -out ~/opt/certificates/apis.company.com.key 2048; \
	openssl ecparam -genkey -name secp384r1 -out ~/opt/certificates/apis.company.com.key; \
	openssl req -new -x509 -sha256 -key ~/opt/certificates/apis.company.com.key -out ~/opt/certificates/apis.company.com.crt -days 365 \
		-subj "/C=US/ST=Tennessee/L=Memphis/O=company Internationl/CN=*local.apis.company.com"

.PHONY: run
#local docker
run: docker
	docker run -p $(HOST_PORT):$(DOCKER_EXPOSED_PORT) -e ENVIRONMENT="local" ${PACKAGE}:latest


.PHONY: run-local
#this is totally cheating, but we rebuild the binary in a subdir of bin (bin/app) so the relative path to the configuration files matches what it is in containers (../../)
#note that you will also need to symlink your certs in (~/opt/certificates) to /opt/certificates so the mount point in containers is consistent.
#this is why run defaults to docker. It's a better testing experience, albeit a little more difficult to debug.
run-local:  
	$(info $(M) building local copy...) @ ## 
	export ENVIRONMENT=local
	$Q mkdir -p bin/app 
	$Q $(GO) build -tags release -ldflags '-X main.GitComHash=$(GITHASH) -X main.BuildStamp=$(DATE)' -o bin/app/application cmd/server/main.go
	$(info $(M) listening on https://localhost:$(DOCKER_EXPOSED_PORT)...) @
	$Q cd bin/app;./application

######################################################
#Dev Tools (linters, build, etc. for development)
######################################################
#lint
.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

#fmt
.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt ./...

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/tests.* test/coverage.*

######################################################
#Testing tools
######################################################
TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt lint
	$(info $(M) clearing $(NAME:%=% )testing cache...) @
	$Q $(GO) clean -testcache ./...
	$(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint | $(GO2XUNIT) 
	$Q mkdir -p test
	$(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml
COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html

.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)