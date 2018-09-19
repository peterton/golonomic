GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=golonomic
EV3_HOST=192.168.158.19

BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
BUILT_AT=`date +%FT%T%z`
BUILT_BY=$(USER)
BUILT_ON=`hostname`
LDFLAGS="-s -w -X main.branch=$(BRANCH) -X main.commit=$(COMMIT) -X main.builtAt='$(BUILT_AT)' -X main.builtBy=$(BUILT_BY) -X main.builtOn=$(BUILT_ON)"

all: test build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -ldflags $(LDFLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy: build stop
	scp $(BINARY_NAME) robot@$(EV3_HOST):

stop:
	ssh -x robot@$(EV3_HOST) "pkill $(BINARY_NAME) || true"

run: deploy 
	ssh -x robot@$(EV3_HOST) ./$(BINARY_NAME)

run_only: stop
	ssh -x robot@$(EV3_HOST) ./$(BINARY_NAME)
