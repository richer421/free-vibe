VERSION=$(shell git describe --tags --always 2>/dev/null || echo dev)

.PHONY: build
# build
build: build-cli

.PHONY: build-cli
# build freevibe cli
build-cli:
	mkdir -p bin/ && go build -ldflags "-X free-vibe-coding/internal/freevibe/cmd.Version=$(VERSION)" -o ./bin/freevibe ./cmd/freevibe

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
