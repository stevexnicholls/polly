.PHONY: test help

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test-verbose:
	go run polly/main.go process --contentDir="test/content" --pollyDir="test/resources/polly" --layoutsDir="test/layouts" --verbose

test:
	go run polly/main.go process --contentDir="test/content" --pollyDir="test/resources/polly" --layoutsDir="test/layouts"