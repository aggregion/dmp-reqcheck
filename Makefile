bin ?= bin
name ?= dmp-reqcheck

linkopts ?= -extldflags=-static -w -s
gobin := ${GOPATH}/bin
target = $(bin)/$(name)
gosrc := $(shell find . -type f -name '*.go' -not -path './vendor/*')
lintfolders := cmd internal pkg
git_commit := $(shell git log -n 1 --pretty=format:'%h')
git_branch := $(shell git rev-parse --abbrev-ref HEAD | tr -dc 'a-zA-Z0-9_-')
app_version := $(shell cat VERSION)

build: $(target)

$(target): $(gosrc)
	@echo 'Build $(target) $(app_version):$(git_branch):$(git_commit)'
	@mkdir -p '$(bin)'
	@cd cmd && go build -o ../$(target) -tags netgo -trimpath -ldflags \
		"-X github.com/aggregion/dmp-reqcheck/internal.GitCommit=$(git_commit) -X github.com/aggregion/dmp-reqcheck/internal.GitBranch=$(git_branch) -X github.com/aggregion/dmp-reqcheck/internal.AppVersion=$(app_version) $(linkopts)"

$(go_bin)/golint:
	@echo "Installing golint"
	@cd /tmp && go get -v -u golang.org/x/lint/golint

lint: $(gobin)/golint
	@echo "Running golint"
	@$(foreach dir,$(lintfolders),golint -set_exit_status $(dir)/...;)
	@echo "Done"

fmt:
	@echo "Formating all go code"
	@$(foreach dir,$(lintfolders),gofmt -s -w $(dir);)
	@echo "Done"

run:
	@echo 'Run $(name)'
	$(bin)/$(name)

test:
	@echo 'Test $(name)'
	@go test -cpu 2 -v -parallel 2 -cover -race ./internal/... ./pkg/...

clean:
	@echo 'clean'
	rm -f $(target)

.PHONY: fmt vendor run test clean
