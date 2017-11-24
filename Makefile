GOPATH = $(CURDIR)/GOPATH
IMPORT_PATH = github.com/pivotal-cf/p-spring-cloud-services-ci-laundromat

export GOPATH

.PHONY: all
all: build

build: setup
	cd $(CURDIR)/GOPATH/src/$(IMPORT_PATH) && go build -i -o bin/laundromat main.go
	cd $(CURDIR)/GOPATH/src/$(IMPORT_PATH) && find vendor -type d -maxdepth 1 -mindepth 1 -exec cp -R "{}" "GOPATH/src" \;

GOPATH/.ok:
	mkdir -p "$(dir GOPATH/src/$(IMPORT_PATH))"
	ln -s ../../../.. "GOPATH/src/$(IMPORT_PATH)"
	mkdir -p GOPATH/test GOPATH/cover
	mkdir -p bin
	ln -s ../bin GOPATH/bin
	touch $@

setup: GOPATH/.ok
	test -s vendor/vendor.json || (echo "vendor.json is missing from vendor folder"; exit 1)
	go get -u github.com/kardianos/govendor
	cd $(CURDIR)/GOPATH/src/$(IMPORT_PATH) && govendor sync

clean:
	rm -rf bin GOPATH