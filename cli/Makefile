# Template Starter CLI Makefile

# 变量定义
BINARY_NAME=template-cli
VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-ldflags "-X github.com/ciclebyte/template_starter/cli/cmd.Version=$(VERSION) -X github.com/ciclebyte/template_starter/cli/cmd.GitCommit=$(GIT_COMMIT) -X github.com/ciclebyte/template_starter/cli/cmd.BuildTime=$(BUILD_TIME)"

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "Template Starter CLI Makefile"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## 构建二进制文件
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

.PHONY: build-all
build-all: ## 构建所有平台的二进制文件
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 main.go

.PHONY: install
install: build ## 安装到本地
	@echo "Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

.PHONY: clean
clean: ## 清理构建文件
	@echo "Cleaning..."
	rm -rf bin/
	go clean

.PHONY: test
test: ## 运行测试
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## 运行代码检查
	@echo "Running linters..."
	golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

.PHONY: mod-tidy
mod-tidy: ## 整理依赖
	@echo "Tidying modules..."
	go mod tidy

.PHONY: mod-download
mod-download: ## 下载依赖
	@echo "Downloading modules..."
	go mod download

.PHONY: dev
dev: fmt mod-tidy build ## 开发构建（格式化+整理依赖+构建）

.PHONY: release
release: clean fmt mod-tidy test build-all ## 发布构建（完整流程）

.PHONY: run
run: ## 运行程序（开发模式）
	go run main.go

.PHONY: run-version
run-version: ## 运行版本命令
	go run main.go version

.PHONY: run-help
run-help: ## 运行帮助命令
	go run main.go --help

.PHONY: docker-build
docker-build: ## 构建Docker镜像
	docker build -t template-cli:$(VERSION) .

.PHONY: docker-run
docker-run: ## 运行Docker容器
	docker run --rm -it template-cli:$(VERSION)

# 开发辅助目标
.PHONY: watch
watch: ## 监视文件变化并重新构建
	@echo "Watching for changes..."
	@which fswatch > /dev/null || (echo "fswatch not found, install with: brew install fswatch" && exit 1)
	fswatch -o . | xargs -n1 -I{} make build

.PHONY: init-dev
init-dev: ## 初始化开发环境
	@echo "Initializing development environment..."
	go mod download
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: check
check: fmt lint test ## 运行所有检查（格式化+代码检查+测试）

# 文档生成
.PHONY: docs
docs: ## 生成命令行文档
	@echo "Generating CLI documentation..."
	go run main.go --help > docs/cli-help.txt
	go run main.go create --help > docs/create-help.txt
	go run main.go template --help > docs/template-help.txt
	go run main.go config --help > docs/config-help.txt

# 依赖相关
.PHONY: deps-update
deps-update: ## 更新所有依赖到最新版本
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: deps-check
deps-check: ## 检查依赖安全性
	@echo "Checking dependencies for security issues..."
	go list -json -m all | nancy sleuth

# 性能分析
.PHONY: bench
bench: ## 运行性能测试
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

.PHONY: profile
profile: ## 运行性能分析
	@echo "Running with profiling..."
	go run -cpuprofile=cpu.prof -memprofile=mem.prof main.go version

# Git 相关
.PHONY: git-hooks
git-hooks: ## 安装Git钩子
	@echo "Installing Git hooks..."
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# 发布相关
.PHONY: tag
tag: ## 创建Git标签
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version: " NEW_VERSION; \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION"; \
	echo "Created tag $$NEW_VERSION"

.PHONY: release-github
release-github: build-all ## 创建GitHub发布
	@echo "Creating GitHub release..."
	@which gh > /dev/null || (echo "GitHub CLI not found, install from https://cli.github.com/" && exit 1)
	gh release create $(VERSION) bin/* --title "Release $(VERSION)" --notes "Release $(VERSION)"

# 清理和维护
.PHONY: clean-all
clean-all: clean ## 清理所有生成的文件
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	go clean -cache
	go clean -modcache