${INSTALL_OPTS}.PHONY: install update test test_unit lint fmt_check fmt
.DEFAULT: help

#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#
#  HELP
#
#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

# COLORS
RED    = $(shell printf "\33[31m")
GREEN  = $(shell printf "\33[32m")
WHITE  = $(shell printf "\33[37m")
YELLOW = $(shell printf "\33[33m")
RESET  = $(shell printf "\33[0m")

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'
# A category can be added with @category
HELP_SCRIPT = \
    %help; \
    while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-\%_]+)\s*:.*\#\#(?:@([a-zA-Z\-\%]+))?\s(.*)$$/ }; \
    print "usage: make [target]\n\n"; \
    for (sort keys %help) { \
    print "${WHITE}$$_:${RESET}\n"; \
    for (@{$$help{$$_}}) { \
    $$sep = " " x (32 - length $$_->[0]); \
    print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
    }; \
    print "\n"; }

help: ##prints help
	@perl -e '${HELP_SCRIPT}' ${MAKEFILE_LIST}

#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#
#  DEPS
#
#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

install: ##@deps Install dependencies
	@echo "${YELLOW}Installing dependencies${RESET}"
	go get ${INSTALL_OPTS} gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get ${INSTALL_OPTS}
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get ${INSTALL_OPTS}
	@echo "${GREEN}✔ successfully installed dependencies${RESET}\n"

update: ##@deps Update dependencies
	@${MAKE} install INSTALL_OPTS=-u

#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#
#  TEST
#
#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

test: ##@test Run all test steps
	@echo "${YELLOW}Running all tests${RESET}\n"
	@${MAKE} test_unit
	@${MAKE} lint
	@${MAKE} fmt_check
	@echo "${GREEN}✔ well done!${RESET}\n"

test_unit: ##@test Run lib tests
	@echo "${YELLOW}Running unit tests${RESET}"
	@go test -v ./...
	@echo "${GREEN}✔ unit tests successfully passed${RESET}\n"

lint: ##@test Run linter
	@echo "${YELLOW}Linting${RESET}"
	@gometalinter.v2 ./...
	@echo "${GREEN}✔ linter passed without error${RESET}\n"

fmt_check: ##@test Check formatting
	@echo "${YELLOW}Checking formatting${RESET}"
	@exit `gofmt -l -s -e . | wc -l`
	@echo "${GREEN}✔ code was formatted as expected${RESET}\n"

#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#
#  UTILS
#
#=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

fmt: ##@utils Format code
	@echo "${YELLOW}Formatting code${RESET}"
	@gofmt -l -w -s .
	@go fix ./...
	@echo "${GREEN}✔ code was successfully formatted${RESET}\n"
