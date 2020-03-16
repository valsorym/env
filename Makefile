# DEVELOPMENT
# Note: parameters for a different projects.
# 	REPOSITORY_NAME
# 		Name of the repository in the GOPATH, like <HOST>/<USER>,
# 		for example: github.com/goloop;
#	PACKAGE_NAME
#		Name of the GoLang's pakcage.
REPOSITORY_NAME="github.com/goloop"
PACKAGE_NAME="env"

# Path to the source of the package.
SRC_PATH:=$(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

# Help information.
define MSG_HELP
The `env` it is environment variable management pack.

Commands:
	help
		Show help information
	link
		Create link to the source code in the GOPATH
	unlink
		Remove link to the source code from the GOPATH
	test
		Run package tests
endef

export MSG_HELP

all: help
help:
	@echo "$$MSG_HELP"
link:
	@cd ${GOPATH}/src/ && \
		mkdir -p ${REPOSITORY_NAME} && cd ${REPOSITORY_NAME} && \
		ln -s $(SRC_PATH) ${PACKAGE_NAME} && \
		echo "Linked: ${GOPATH}/src/${REPOSITORY_NAME}/${PACKAGE_NAME}" && \
		ls -l ${PACKAGE_NAME}
unlink:
	@cd ${GOPATH}/src/ && \
		rm -Rf ${REPOSITORY_NAME}/${PACKAGE_NAME} && \
		echo "Unlinked: ${GOPATH}/src/${REPOSITORY_NAME}/${PACKAGE_NAME}"
test:
	@go test github.com/goloop/env
test-cover:
	@go test -cover github.com/goloop/env && \
		go test -coverprofile=/tmp/coverage.out github.com/goloop/env && \
		go tool cover -func=coverage.out && \
		go tool cover -html=coverage.out
