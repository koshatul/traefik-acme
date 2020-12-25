DOCKER_REPO := koshatul/traefik-acme

GO_MATRIX_OS ?= linux darwin
GO_MATRIX_ARCH ?= amd64

APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH ?= $(shell git show -s --format=%h)


GO_DEBUG_ARGS   ?= -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug -X main.commit=$(GIT_HASH) -X main.date=$(APP_DATE) -X main.builtBy=makefiles.dev"
GO_RELEASE_ARGS ?= -v -ldflags "-s -w -X main.version=$(GO_APP_VERSION) -X main.commit=$(GIT_HASH) -X main.date=$(APP_DATE) -X main.builtBy=makefiles.dev"

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
# Tests
######################

TEST_RUNNER := artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/traefik-acme

.DELETE_ON_ERROR: artifacts/test/issue-5/test1.out
REGRESSION_TESTS += artifacts/test/issue-5/test1.out
artifacts/test/issue-5/test1.out: test/issue-5/acme.json $(TEST_RUNNER) $(GO_SOURCE_FILES)
	-@mkdir -p "$(@D)"
	-@$(RM) "$(@D)/cert.pem" "$(@D)/key.pem"
	$(TEST_RUNNER) --acme="$(<)" --certificate-resolver="acme-different" --cert "$(@D)/cert.pem" --key "$(@D)/key.pem" "*.example.com" | tee "$(@)"
	grep "certificate-for-example.com" "$(@D)/cert.pem"
	grep "key-for-example.com" "$(@D)/key.pem"

.DELETE_ON_ERROR: artifacts/test/issue-5/test2.out
REGRESSION_TESTS += artifacts/test/issue-5/test2.out
artifacts/test/issue-5/test2.out: test/issue-5/acme.json $(TEST_RUNNER) $(GO_SOURCE_FILES)
	-@mkdir -p "$(@D)"
	-@$(RM) "$(@D)/cert.pem" "$(@D)/key.pem"
	$(TEST_RUNNER) --acme="$(<)" --certificate-resolver="acme-different" --cert "$(@D)/cert.pem" --key "$(@D)/key.pem" "example.com" | tee "$(@)"
	grep "certificate-for-example.com" "$(@D)/cert.pem"
	grep "key-for-example.com" "$(@D)/key.pem"

.PHONY: regression-tests
regression-tests: $(REGRESSION_TESTS)

ci:: $(REGRESSION_TESTS)

######################
# Linting
######################

MISSPELL := artifacts/bin/misspell
$(MISSPELL):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/client9/misspell/cmd/misspell

GOLINT := artifacts/bin/golint
$(GOLINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) golang.org/x/lint/golint

GOLANGCILINT := artifacts/bin/golangci-lint
$(GOLANGCILINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(MF_PROJECT_ROOT)/$(@D)" v1.33.0

STATICCHECK := artifacts/bin/staticcheck
$(STATICCHECK):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) honnef.co/go/tools/cmd/staticcheck

artifacts/cover/staticheck/unused-graph.txt: $(STATICCHECK) $(GO_SOURCE_FILES)
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	$(STATICCHECK) -debug.unused-graph "$(@)" ./...
	# cat "$(@)"

.PHONY: lint
lint:: $(GOLINT) $(MISSPELL) $(GOLANGCILINT) $(STATICCHECK) artifacts/cover/staticheck/unused-graph.txt
	go vet ./...
	$(GOLINT) -set_exit_status ./...
	$(MISSPELL) -w -error -locale UK ./...
	$(GOLANGCILINT) run --enable-all --disable 'exhaustivestruct,paralleltest' ./...
	$(STATICCHECK) -fail "all,-U1001" ./...

ci:: lint


######################
# Preload Tools
######################

.PHONY: tools
tools: $(MISSPELL) $(GOLINT) $(GOLANGCILINT) $(STATICCHECK)
