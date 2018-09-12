GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=golonomic
EV3_HOST=192.168.200.73

all: test build
	
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	scp $(BINARY_NAME) robot@$(EV3_HOST):
	ssh robot@$(EV3_HOST) ./$(BINARY_NAME)
