.PHONY: build 
.PHONY: push 
.PHONY: run 

GO=go
VERSION_FILE=./cmd/version
CURRENT_VERSION=`cat $(VERSION_FILE)`

BUILD_TAGS=debugcharts

build:
	GOOS=linux GOARCH=amd64 go build -tags='$(BUILD_TAGS)' -o ./bin/ae-copilot  ./cmd
	docker build -f ./cmd/Dockerfile.local -t ae-copilot:latest .
 
push:
	@echo 'Current ae-copilot version: '$(COUNTS_CURRENT_VERSION);\
	read -p "New version? " input;\
	docker tag ae-copilot:latest liveramp-cn-north-1.jcr.service.jdcloud.com/ae-copilot:$${input};\
	docker push liveramp-cn-north-1.jcr.service.jdcloud.com/ae-copilot:$${input};\
	echo $${input} > $(COUNTS_VERSION_FILE);\
	git commit $(COUNTS_VERSION_FILE) -m "ci: build new version $${input}"
 

run:
	@echo 'Current ae-copilot version: '$(COUNTS_CURRENT_VERSION)
	cd ./cmd; go run main.go --addr :5001
 