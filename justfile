# Set environment variables
export ENV := env_var_or_default("ENV", "dev")
export PACKAGES := "./..."

# Load .env file if it exists
set dotenv-load

# Use 1Password only if SHORTCUT_API_TOKEN is a secret reference (op://...)
op := if env("SHORTCUT_API_TOKEN", "") =~ '^op://' { "op run --no-masking --" } else { "" }

# Default recipe to display help
default:
    @just --list

# Initialize requirements
init:
    go tool swagger generate client -f https://developer.shortcut.com/api/rest/v3/shortcut.swagger.json --target pkg/shortcut/gen/
    go mod download

# Build in production mode
build:
    goreleaser release --clean

# Build in snapshot mode
build-snapshots:
    goreleaser release --clean --snapshot

# Lint the project
lint:
    golangci-lint run $PACKAGES

# Execute application code with optional arguments
dev *args:
    {{op}} go run cmd/cli/main.go {{args}}

# Launch tests with coverage
tests:
    {{op}} go test -cover $PACKAGES

# Run tests with verbose output
test-verbose:
    {{op}} go test -v -cover $PACKAGES

# Run tests and generate coverage report
test-coverage:
    {{op}} go test -coverprofile=coverage.out $PACKAGES
    go tool cover -html=coverage.out
