#!/bin/bash
PLUGIN_DIR=plugins
BUILD_DIR=bin
PLUGIN_DIR_OUT=${BUILD_DIR}/plugins
ENGINE_DIR_OUT=${BUILD_DIR}/
PLUGINS_LIST=plugin1 plugin2


all: lint build-plugins build

build:
	@go build -o ${BUILD_DIR}/

build-plugins:
	@mkdir -p ${PLUGIN_DIR_OUT}
	@for plugin in ${PLUGINS_LIST}; do \
		echo "Building $$plugin ..."; \
		go build -buildmode=plugin -o ${PLUGIN_DIR_OUT}/$$plugin.so ./plugins/src/$$plugin ; \
	done
	

clean:
	rm -rf ${PLUGIN_DIR_OUT}


lint:
	golangci-lint run


toolkit:
	brew install  golang golangci-lint
