.PHONY: default build build-all release-local install lint test testacc clean

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=datadome.co
NAMESPACE=app
NAME=datadome
BINARY=terraform-provider-${NAME}
# the version of the local binary that will be generated
VERSION=0.0.1
OS_ARCH=`uname -s | tr A-Z a-z`_`uname -m | tr A-Z a-z`

default: install

build:
	go build -o ${BINARY}

build-all:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./dist/${BINARY}_${VERSION}_darwin_arm64
	GOOS=freebsd GOARCH=386 go build -o ./dist/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./dist/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./dist/${BINARY}_${VERSION}_freebsd_arm
	GOOS=freebsd GOARCH=arm64 go build -o ./dist/${BINARY}_${VERSION}_freebsd_arm64
	GOOS=linux GOARCH=386 go build -o ./dist/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./dist/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./dist/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=arm64 go build -o ./dist/${BINARY}_${VERSION}_linux_arm64
	GOOS=openbsd GOARCH=386 go build -o ./dist/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./dist/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=windows GOARCH=386 go build -o ./dist/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./dist/${BINARY}_${VERSION}_windows_amd64
	GOOS=windows GOARCH=arm go build -o ./dist/${BINARY}_${VERSION}_windows_arm
	GOOS=windows GOARCH=arm64 go build -o ./dist/${BINARY}_${VERSION}_windows_arm64

release-local:
	goreleaser release --snapshot --rm-dist --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

lint:
ifeq (, $(shell which golangci-lint))
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
endif
	golangci-lint run ./main.go
	golangci-lint run ./datadome/*.go
	golangci-lint run ./datadome-client-go/*.go

test: 
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

clean:
	rm -rf ./dist/
