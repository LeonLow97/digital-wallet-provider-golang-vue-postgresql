GO 				:= $(shell which go)
PACKAGES 		:= $(shell $(GO) list ./... | grep -v /vendor)
COVERAGE_FILE 	:= coverage.out

# Build target
.PHONY: build
build: format
	docker-compose down
	docker-compose up --build -d
	@echo "Client available at http://localhost:3000"
	@echo "Server API available at http://localhost:8080"

# Up target
.PHONY: up
up: format
	docker-compose up --build -d
	@echo "Client available at http://localhost:3000"
	@echo "Server API available at http://localhost:8080"

# Down target
.PHONY: down
down: format
	docker-compose down

# Stop target
.PHONY: stop
stop: format
	docker-compose stop

# Unit test target
.PHONY: unit
unit: format
	$(GO) test -race -v -count=1 $(PACKAGES) -coverprofile=$(COVERAGE_FILE)

# Coverage HTML target
.PHONY: coverage-html
coverage-html: unit
	$(GO) tool cover -html=$(COVERAGE_FILE)

# Formatting target
.PHONY: format
format: 
	gofmt -w .