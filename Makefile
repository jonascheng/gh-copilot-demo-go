.DEFAULT_GOAL := help


.PHONY: setup
setup: ## setup python package
	go mod tidy


.PHONY: run
run: setup ## runs streamlit application
	go run app.go


.PHONY: help
help: ## prints this help message
	@echo "Usage: \n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
