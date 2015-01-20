PACKAGES=github.com/patkaehuaea/server
GODOC_PORT=:6060

all: fmt install

install:
	GOPATH=$(GOPATH) go install $(PACKAGES)

test:
	GOPATH=$(GOPATH) go test $(TEST_PACKAGES)

fmt:
	GOPATH=$(GOPATH) go fmt $(PACKAGES)

doc:
	GOPATH=$(GOPATH) godoc -v --http=$(GODOC_PORT) --index=true

clean:
	go clean