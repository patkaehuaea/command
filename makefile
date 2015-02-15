PACKAGES=./...
ARTIFACT_DIRS=$(GOPATH)/bin $(GOPATH)/out $(GOPATH)/pkg
GO_PARENT_DIR=$(HOME)
ZIP_DEST_DIR=$(HOME)
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
	GOPATH=$(GOPATH) go clean

delartifacts:
	@for i in $(ARTIFACT_DIRS); do \
        echo "Deleteing files in $$i..."; \
        GOPATH=$(GOPATH) /bin/rm -rf $$i/*; \
    done

zip:
	(cd $(GO_PARENT_DIR) ; /usr/bin/zip -r $(ZIP_DEST_DIR)/go.zip go)

release: clean delartifacts zip
