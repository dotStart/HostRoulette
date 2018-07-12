APPLICATION_VERSION := 0.1.0
APPLICATION_COMMIT_HASH := `git log -1 --pretty=format:"%H"`

LDFLAGS :=-ldflags "-X github.com/dotStart/HostRoulette/command.version=${APPLICATION_VERSION} -X github.com/dotStart/HostRoulette/command.commitHash=${APPLICATION_COMMIT_HASH}"

GIT := $(shell command -v git 2> /dev/null)
DEP := $(shell command -v dep 2> /dev/null)
GO := $(shell command -v go 2> /dev/null)
NODE := $(shell command -v node 2> /dev/null)
NPM := $(shell command -v npm 2> /dev/null)

ifdef SYSTEMROOT
EXT=".exe"
else
EXT=""
endif

all: check-env print-config install-dependencies compile-ui compile-server

check-env:
	@echo "==> Checking prerequisites"
	@echo -n "Checking for git ... "
ifndef GIT
	@echo "Not found"
	$(error "git is unavailable")
endif
	@echo $(GIT)
	@echo -n "Checking for go ... "
ifndef GO
	@echo "Not Found"
	$(error "go is unavailable")
endif
	@echo $(GO)
	@echo -n "Checking for dep ... "
ifndef DEP
	@echo "Not Found"
	$(error "dep is unavailable")
endif
	@echo $(DEP)
	@echo -n "Checking for node ... "
ifndef NODE
	@echo "Not Found"
	$(error "node is unavailable")
endif
	@echo $(NODE)
	@echo -n "Checking for npm ... "
ifndef NPM
	@echo "Not Found"
	$(error "npm is unavailable")
endif
	@echo $(NPM)
	@echo ""

print-config:
	@echo "==> Build Configuration"
	@echo ""
	@echo "     Version: ${APPLICATION_VERSION}"
	@echo "  Commit SHA: ${APPLICATION_COMMIT_HASH}"
	@echo ""

clean:
	@echo "==> Clearing previous build data"
	@rm -rf build/ || true
	@$(GO) clean -cache

install-dependencies:
	@echo "==> Installing dependencies"
	@$(GO) get -u github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	@$(GO) get -u github.com/jteeuwen/go-bindata/go-bindata
	@$(DEP) ensure -v
	@echo ""

.ONESHELL:
compile-ui:
	@echo "==> Compiling UI"
	cd ui
	"$(NPM)" install
	"$(NPM)" run build

compile-server:
	@echo "==> Building"
	$(GO) generate -v github.com/dotStart/HostRoulette/ui
	$(GO) build -v ${LDFLAGS} -o build/HostRoulette${EXT}

.PHONY: all
