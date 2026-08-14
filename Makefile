.PHONY: run run-dev server # ensure these will be run as commands

MAKEFLAGS += --no-print-directory # ignores extra messages from make

ENV ?= prod # environment (production or dev)
ID ?= 197 # pokemon ID (currently defaults to Umbreon)
WRITE_DB ?= 0 # whether or not to export to database

run:
	go mod tidy
	P7D_ENV=$(ENV) P7D_WRITE_DB=$(WRITE_DB) go run -buildvcs=true . $(ID)

#Note: $(MAKE) expands to the making command, or 'make' in bash
#So this just translates to make run ENV=dev

run-dev:
	$(MAKE) run ENV=dev 

# run with database export
super-run-dev:
	$(MAKE) run-dev WRITE_DB=1

server:
	go run ./server/.