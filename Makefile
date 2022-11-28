.PHONY: init

init:
	@echo "Initialize development environment..."
	@mkdir -p .dev
	@go mod tidy

.PHONY: deps

deps:
	@echo "Install dependencies..."
	@go mod tidy

.PHONY: clean

clean:
	@echo "Clean development files..."
	@rm -rf .dev
