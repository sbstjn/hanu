COVERAGE_FILE ?= c.out

test:
	@ ginkgo -cover -coverprofile=$(COVERAGE_FILE) $(RACE) ./... 

lint:
	@ golint ./..

tool:
	@ go tool cover -$(MODE)=$(COVERAGE_FILE)

race: RACE=-race
race: test

func: MODE=func
func: test tool

html: MODE=html
html: test tool

.PHONY: test lint tool race func html