GOPATH?=$(shell go env GOPATH)
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin
GOSRC:=$(FIRST_GOPATH)/src

VERSION?=$(shell cat VERSION || echo "1.0.0")
BUILD=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git tag --points-at HEAD)
BUILD_DATE=$(shell date -u +%Y-%m-%d-%H:%M)

LDFLAGS=-ldflags " -X main.Version=${VERSION} \
	-X ${GOLIBS_APP_PKG}.Build=${BUILD} \
	-X ${GOLIBS_APP_PKG}.BuildTag=${GIT_TAG} \
	-X ${GOLIBS_APP_PKG}.BuildDate=${BUILD_DATE}"
