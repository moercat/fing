# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOMOD) tidy
GOVET=$(GOCMD) vet
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Install dependencies
tidy:
	$(GOTIDY)

# Run tests
test:
	$(GOTEST) -v ./...

# Run vet to find suspicious constructs
vet:
	$(GOVET) ./...

# Run all checks
check: vet test

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -a -installsuffix cgo .

# Run docker build
docker-build:
	docker build -t fing-app .

# Run docker-compose
docker-up:
	docker-compose up -d

# Stop docker-compose
docker-down:
	docker-compose down

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

.PHONY: build run test vet check build-prod docker-build docker-up docker-down clean tidy