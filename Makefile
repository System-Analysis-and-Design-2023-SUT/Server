##
# Base Makefile with ability to list targets
#
.PHONY: help about list targets

ME := $(realpath $(firstword $(MAKEFILE_LIST)))

# Contains trailing '/'
#
PWD := $(dir $(ME))

.DEFAULT_GOAL := help

##
# help
# Displays a (hopefully) useful help screen to the user
#
# NOTE: Keep 'help' as first target in case .DEFAULT_GOAL is not honored
#
help: about list ## This help screen
about:
	@echo
	@echo "Makefile to help manage System Analysis and Design Project Server"

##
# list
# Displays a list of targets, using '##' comment as target description
#
# NOTE: ONLY targets with ## comments are shown
#
list: targets ## see 'targets'
targets: ## Lists targets
	@echo
	@echo  "Make targets:"
	@echo
	@cat $(ME) | \
	sed -n -E 's/^([^.][^: ]+)\s*:(([^=#]*##\s*(.*[^[:space:]])\s*)|[^=].*)$$/    \1	\4/p' | \
	sort -u | \
	expand -t15
	@echo

test:
	go test ./... -coverprofile cover.out

build:
	go env -w GO111MODULE="on"
	go build -a -o bin/app cmd/main.go

run:
	./bin/app
