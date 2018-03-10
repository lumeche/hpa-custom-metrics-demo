REGISTRY?=lumeche
SS_IMAGE?=scalable-server
CM_IMAGE?=custom-metrics
TEMP_DIR:=$(shell mktemp -d)
VERSION?=latest

all: build_SS build_CM push_CM push_SS
build_SS: 
	docker run -it -v $(TEMP_DIR):/build -v $(shell pwd)/scalable_server:/go/src/github.com/lumeche/scalable_server -e GOARCH=$(ARCH) golang:1.9 /bin/bash -c "\
		CGO_ENABLED=0 go build -a -tags netgo -o /build/scalable_server github.com/lumeche/scalable_server"

	
	cp scalable_server/Dockerfile $(TEMP_DIR)
	docker build -t $(REGISTRY)/$(SS_IMAGE):$(VERSION) $(TEMP_DIR)
	rm -rf $(TEMP_DIR)
	

build_CM: 
	docker run -it -v $(TEMP_DIR):/build -v $(shell pwd)/metrics_server:/go/src/github.com/lumeche/metrics_server -e GOARCH=$(ARCH) golang:1.9 /bin/bash -c "\
		CGO_ENABLED=0 go build -a -tags netgo -o /build/metrics_server github.com/lumeche/metrics_server"
	cp metrics_server/Dockerfile $(TEMP_DIR)
	docker build -t $(REGISTRY)/$(CM_IMAGE):$(VERSION) $(TEMP_DIR)
	rm -rf $(TEMP_DIR)


push_SS:
	docker push $(REGISTRY)/$(SS_IMAGE):$(VERSION) 
push_CM:
	docker push $(REGISTRY)/$(CM_IMAGE):$(VERSION) 
