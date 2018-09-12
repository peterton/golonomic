GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=golonomic
EV3_HOST=192.168.200.73

all: test build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	scp $(BINARY_NAME) $(EV3_HOST):
	ssh $(EV3_HOST) ./$(BINARY_NAME)
deps:
	$(GOGET) github.com/ev3go/ev3dev
	$(GOGET) gonum.org/v1/gonum/mat
