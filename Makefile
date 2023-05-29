.PHONY: build

PWD := $(shell pwd)

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

PROTO_DIR := /cmd/monitor/proto

#initFirst:
#	export GOSUMDB=off

proto:
	# adapter for grpc 1.26.0
	# go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
	@echo 'generating go from proto files'
	@echo $(PWD)$(PROTO_DIR)
    ifeq ($(file), "")
		@protoc -I$(PWD)$(PROTO_FILE) --go_out=plugins=grpc:$(PWD)$(PROTO_DIR) $(PWD)$(PROTO_DIR)/*.proto
	else
  		@protoc -I$(PWD)$(PROTO_FILE) --go_out=plugins=grpc:$(PWD)$(PROTO_DIR) $(PWD)$(PROTO_DIR)/$(file).proto
    endif

db:
	xorm-mac reverse mysql 'other_rw:DA65d357D8dd4666bf4fAbfD6624f139@(10.90.29.171:6306)/xesflow?charset=utf8mb4' ./goxorm models

build:
	go build -ldflags $(LD_FLAGS) -gcflags "-N" -o ./bin/main ./cmd/monitor

monitor:
	docker run -d --name promethu -p 9090:9090 -v /Users/huaxinrui/docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
	&& docker run -d -p 3000:3000 --name=grafana grafana/grafana
	&& docker run -d --name pushgateway -p 9091:9091 prom/pushgateway