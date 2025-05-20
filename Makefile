# This Makefile defines common tasks for a Go project, including linting the code,
# running tests, and generating mocks using mockery.

# Target: lint
# Description: Run the Go linter using golangci-lint.
# This checks the code for stylistic issues, potential bugs, and other improvements.
lint: ## Run linter
	golangci-lint run

# Target: test
# Description: Run all tests in the project.
# The -v flag makes the test output verbose, showing detailed information about each test.
test: ## Run tests
	go test -v ./...

# Target: bump
# Description: Bump version, generate git tag, and push it to the repository.
bump: semver-cli ## Bump version, generate git tag and push it to repository
	@printf "\033[36m%s\033[0m\n" "Bumping version..."
	@git fetch --all --tags
	@semver-cli tags bump -t patch -p v

# Target: mocks
# Purpose: Regenerates mocks for interfaces in the project using mockery. Installs mockery if not already available.
mocks: ## Generate mocks
ifeq (, $(shell which mockery))
	go install github.com/vektra/mockery/v2@latest
endif
	mockery