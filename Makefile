.DEFAULT_GOAL := examples-info

.PHONY: check
check:
	build/check-go

.PHONY: gen
gen:
	go generate ./...

.PHONY: examples-basic
examples-basic: check gen
	go run ./examples/basic/...

.PHONY: examples-info
examples-info: check gen
	go run ./examples/info/...

.PHONY: examples-redshift
examples-redshift: check gen
	go run ./examples/redshift/...
