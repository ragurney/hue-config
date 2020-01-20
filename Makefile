# Required Env
GO111MODULE=on

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# General parameters
CURR_DIR=$(shell pwd)

# SAM parameters
SAM_BUILD_DIR=.aws-sam/build/

# Lambda handlers
alexa: ./lambdas/alexa/main.go
	$(GOBUILD) -o $(SAM_BUILD_DIR)alexa ./lambdas/alexa

authentication: ./lambdas/authentication/main.go
	$(GOBUILD) -o $(SAM_BUILD_DIR)authentication ./lambdas/authentication

.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) alexa
	GOOS=linux GOARCH=amd64 $(MAKE) authentication

.PHONY: clean
clean:
	rm -rfv $(SAM_BUILD_DIR)
	mkdir $(SAM_BUILD_DIR)

.PHONY: build
build: clean lambda
	cp template.yaml $(SAM_BUILD_DIR)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: run
run: build
	sam local start-api

.PHONY: deploy
deploy: test build
	sam deploy --guided

.PHONY: all
all: test build