DOCKER_REPO := koshatul/traefik-acme

GO_MATRIX_OS ?= darwin linux
GO_MATRIX_ARCH ?= amd64

APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH ?= $(shell git show -s --format=%h)

GO_DEBUG_ARGS   ?= -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE)"
GO_RELEASE_ARGS ?= -v -ldflags "-X main.version=$(GO_APP_VERSION) -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE) -s -w"

_GO_GTE_1_14 := $(shell expr `go version | cut -d' ' -f 3 | tr -d 'a-z' | cut -d'.' -f2` \>= 14)
ifeq "$(_GO_GTE_1_14)" "1"
_MODFILEARG := -modfile tools.mod
endif

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/pkg/docker/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

.PHONY: install
install: artifacts/build/release/$(GOHOSTOS)/$(GOHOSTARCH)/traefik-acme
	install "$(<)" /usr/local/bin/

.PHONY: run
run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/traefik-acme
	$< $(RUN_ARGS)

.PHONY: upx
upx: $(patsubst artifacts/build/%,artifacts/upx/%.upx,$(_GO_RELEASE_TARGETS_ALL))

artifacts/upx/%.upx: artifacts/build/%
	-@mkdir -p "$(@D)"
	-$(RM) -f "$(@)"
	upx -o "$@" "$<"

.PHONY: docker-run
docker-run: docker
	docker run -ti --rm $(DOCKER_REPO):$(DOCKER_TAGS) $(RUN_ARGS)

# .PHONY: run
# run: artifacts/build/debug/$(GOOS)/$(GOARCH)/traefik-acme
# 	$< $(RUN_ARGS)

######################
# Linting
######################

MISSPELL := artifacts/misspell/bin/misspell
$(MISSPELL):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/client9/misspell/cmd/misspell

GOLINT := artifacts/golint/bin/golint
$(GOLINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) golang.org/x/lint/golint

GOLANGCILINT := artifacts/golangci-lint/bin/golangci-lint
$(GOLANGCILINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(MF_PROJECT_ROOT)/$(@D)" v1.23.8

STATICCHECK := artifacts/staticcheck/bin/staticcheck
$(STATICCHECK):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) honnef.co/go/tools/cmd/staticcheck

.PHONY: lint
lint:: $(MISSPELL) $(GOLINT) $(GOLANGCILINT) $(STATICCHECK)
	go vet ./...
	$(GOLINT) -set_exit_status ./...
	$(MISSPELL) -w -error -locale UK ./...
	$(GOLANGCILINT) run --enable-all --disable=dupl,gomnd ./...
	$(STATICCHECK) -checks "all" -fail "all,-U1001" -unused.whole-program ./...


######################
# Preload Tools
######################

.PHONY: tools
tools: $(MISSPELL) $(GOLINT) $(GOLANGCILINT) $(STATICCHECK)
