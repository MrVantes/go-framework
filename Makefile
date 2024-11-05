# List of projects to build
PROJECTS := application

# Database configuration
DATABASE_URL := postgres://postgres:12adib%40was@127.0.0.1/nota_db?sslmode=disable

# Tool versions and paths
GOLANGCI_LINT_VERSION := v1.50.1
LINTER := bin/golangci-lint
MIGRATE := migrate
SQLC := sqlc

# Build targets
all: build

clean: TARGET=clean
clean: default

build: TARGET=all
build: PROJECTS:=$(PROJECTS)
build: default

release: TARGET=release
release: default 

docker-build: TARGET=docker-build
docker-build: default

docker-push: TARGET=docker-push
docker-push: default

# Linting
$(LINTER):
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION)

lint: TARGET=lint
lint: PROJECTS:=$(PROJECTS)
lint: $(LINTER) default

# Repository setup commands
.PHONY: setup-repo check-github-user replace-module-name

check-github-user:
	@if [ "$(GITHUB_USER)" = "GITHUB_USER_NOT_SET" ]; then \
		echo "Error: GitHub username not set. Please set it using:"; \
		echo "make setup-repo GITHUB_USER=your-github-username"; \
		exit 1; \
	fi

replace-module-name: check-github-user
	@echo "Updating module name from $(DEFAULT_MODULE) to $(NEW_MODULE)"
	@# Update go.mod files
	@find . -name "go.mod" -exec sed -i 's|$(DEFAULT_MODULE)|$(NEW_MODULE)|g' {} +
	@# Update all Go imports
	@find . -type f -name "*.go" -exec sed -i 's|$(DEFAULT_MODULE)|$(NEW_MODULE)|g' {} +
	@# Update Makefiles
	@find . -name "Makefile" -exec sed -i 's|$(DEFAULT_MODULE)|$(NEW_MODULE)|g' {} +
	@# Update any Docker/deployment files
	@find . -type f -path "*/deploy/*" -exec sed -i 's|$(DEFAULT_MODULE)|$(NEW_MODULE)|g' {} +

setup-repo: replace-module-name
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "Cleaning build artifacts..."
	@$(MAKE) clean
	@echo "Building project..."
	@$(MAKE) build
	@echo "\nRepository setup complete for $(GITHUB_USER)"
	@echo "Next steps:"
	@echo "1. Review the changes (git status)"
	@echo "2. Commit the changes (git commit)"
	@echo "3. Update your remote origin:"
	@echo "   git remote set-url origin https://github.com/$(GITHUB_USER)/$(PROJECTS).git"
	@echo "4. Push to your repository:"
	@echo "   git push -u origin main"

# Default build process
default:
	@for PRJ in $(PROJECTS); do \
		echo "--- $$PRJ: $(TARGET) ---"; \
		$(MAKE) $(TARGET) -C $$PRJ || exit 1; \
	done

# Database migrations
.PHONY: migrate-up migrate-down generate copy-migrations

migrate-up:
	$(MIGRATE) -path ${PROJECTS}/db/migrations -database $(DATABASE_URL) up

migrate-down:
	$(MIGRATE) -path ${PROJECTS}/db/migrations -database $(DATABASE_URL) down --all

generate:
	$(SQLC) generate -f ${PROJECTS}/sqlc.yaml

# Development commands
.PHONY: run format test

run:
	make clean && make && \
	DATABASE_URL=$(DATABASE_URL) ./${PROJECTS}/build/bin/${PROJECTS}-server

format:
	find . -name \*.go -exec goimports -w {} \;

test:
	go test ./... -v

copy-migrations:
	cp -r ${PROJECTS}/db/migrations deploy/compose

# Help command to show available commands
help:
	@echo "Available commands:"
	@echo "  setup-repo GITHUB_USER=<username>  - Setup repository for your GitHub account"
	@echo "  build                              - Build all projects"
	@echo "  clean                              - Clean build artifacts"
	@echo "  test                               - Run tests"
	@echo "  lint                               - Run linter"
	@echo "  format                             - Format code"
	@echo "  run                                - Run the server"
	@echo "  migrate-up                         - Run database migrations"
	@echo "  migrate-down                       - Revert database migrations"
	@echo "  generate                           - Generate SQLC code"
	@echo "  docker-build                       - Build Docker images"
	@echo "  docker-push                        - Push Docker images"
	@echo "  copy-migrations                    - Copy migrations to deploy directory"

.PHONY: all $(PROJECTS) clean build docker-build docker-push release test lint default help