GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=golonomic
EV3_HOST=192.168.221.132

COMMIT=`git rev-parse --short HEAD`
BUILT_AT=`date +%FT%T%z`
BUILT_BY=$(USER)
BUILT_ON=`hostname`
LDFLAGS="-s -w -X main.commit=$(COMMIT) -X main.builtAt='$(BUILT_AT)' -X main.builtBy=$(BUILT_BY) -X main.builtOn=$(BUILT_ON)"

all: test build
	
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -ldflags $(LDFLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy: build
	ssh -x robot@$(EV3_HOST) "pkill $(BINARY_NAME) || true"
	scp $(BINARY_NAME) robot@$(EV3_HOST):

run: deploy
	ssh -x robot@$(EV3_HOST) ./$(BINARY_NAME)

run_only:
	ssh -x robot@$(EV3_HOST) "pkill $(BINARY_NAME) || true"
	ssh -x robot@$(EV3_HOST) ./$(BINARY_NAME)
