#!/usr/bin/bash

.PHONY: all

DEPLOYMENT_TOKEN    ?= ***REMOVED***
SNAP_CACHE_DIR      ?= $(HOME)/.appengine-go
APPENGINE_PATH      ?= $(SNAP_CACHE_DIR)/go_appengine
GOPATH              ?= $(APPENGINE_PATH)/gopath
GOBIN               := $(GOPATH)/bin
GIGO                := $(GOBIN)/gigo
GO                  := $(APPENGINE_PATH)/goapp
APPCFG              := $(APPENGINE_PATH)/appcfg.py
APP_IMPORT_PATH = github.com/bearchinc/trails-api
APP_PATH        := $(GOPATH)/src/$(APP_IMPORT_PATH)

env:
	@echo "Using the following build parameters:"
	@echo ""
	@echo "DEPLOYMENT_TOKEN=******"
	@echo "SNAP_CACHE_DIR=$(SNAP_CACHE_DIR)"
	@echo "APPENGINE_PATH=$(APPENGINE_PATH)"
	@echo "GOPATH=$(GOPATH)"
	@echo "GOBIN=$(GOBIN)"
	@echo "GIGO=$(GIGO)"
	@echo "GO=$(GO)"
	@echo "APPCFG=$(APPCFG)"
	@echo "APP_IMPORT_PATH=$(APP_IMPORT_PATH)"
	@echo "APP_PATH=$(APP_PATH)"
	@echo ""

setup: env
	@./bin/download-appengine.sh $(SNAP_CACHE_DIR)
	@./bin/resolve-dependencies.sh $(GOPATH) $(GO)

test: setup update
	$(GO) test ./... -v -run=$(grep)

build: setup update
	$(GO) build ./...

update: setup
	$(GO) get github.com/golang/protobuf/proto
	$(GO) get google.golang.org/appengine/urlfetch
	$(GO) get google.golang.org/grpc
	$(GO) get -u golang.org/x/oauth2/google
	$(GO) get golang.org/x/net/context
	$(GO) get -u github.com/drborges/rivers
	$(GO) get github.com/go-martini/martini
	$(GO) get github.com/martini-contrib/render
	$(GO) get github.com/martini-contrib/binding
	$(GO) get google.golang.org/cloud/compute/metadata
	$(GO) get github.com/dchest/authcookie
	$(GO) get github.com/drborges/appx

	rm -rf $(GOPATH)/src/github.com/stacktic/dropbox 2> /dev/null
	@git clone https://github.com/BearchInc/dropbox.git $(GOPATH)/src/github.com/stacktic/dropbox
	@git clone https://github.com/BearchInc/geocoder.git $(GOPATH)/src/github.com/drborges/geocoder

delete-branches:
	git branch | grep -v master | xargs -I {} git branch -D {}

serve:
	goapp serve --host 0.0.0.0 app/app.yaml

serve_fresh:
	goapp serve -clear_datastore --host 0.0.0.0 app/app.yaml

put-app-in-go-path:
	mkdir -p $(GOPATH)/src/github.com/bearchinc
	ln -sf `pwd` $(APP_PATH)

ci-build-info:
	./bin/build-info.sh ./app/info/build.json $(SNAP_STAGE_NAME) $(SNAP_PIPELINE_COUNTER) $(SNAP_BRANCH) $(SNAP_COMMIT)

ci-deploy-staging: setup update put-app-in-go-path ci-build-info
	$(APPCFG) --oauth2_refresh_token=$(DEPLOYMENT_TOKEN) -A staging-api-getunseen update_indexes app
	$(APPCFG) --oauth2_refresh_token=$(DEPLOYMENT_TOKEN) -A staging-api-getunseen update app

# Deployment tasks
deploy:
	goapp deploy -oauth app

rollback-deploy:
	appcfg.py --oauth2_refresh_token=***REMOVED*** rollback app

# Handy curls
http-get:
	curl -XGET https://trails-dot-staging-api-getunseen.appspot.com$(path)

http-post:
	curl -XPOST https://trails-dot-staging-api-getunseen.appspot.com$(path) \
		-H "Content-type: application/json" \
		-d '$(json)'
