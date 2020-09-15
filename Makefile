.PHONY: all
all: ./bin/yt-relay

bin/%: $(shell find . -type f -name '*.go')
	@mkdir -p $(dir $@)
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -tags=jsoniter -o $@ ./cmd/$(@F)


.PHONY: clean
clean:
	rm -rf bin
