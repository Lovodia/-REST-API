# Имя бинарника
BINARY_NAME=calculator-api
BUILD_DIR=build

# Путь к Swagger
SWAGGER_DIR=./docs

# Цель по умолчанию
all: build

## Сборка бинарника
build:
	@echo "Building the binary..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

## Генерация Swagger документации
swagger:
	@echo "Generating Swagger docs..."
	swag init --output $(SWAGGER_DIR) --parseDependency --parseInternal

## Запуск тестов
test:
	@echo "Running tests..."
	go test -v ./...

## Запуск локально
run:
	@echo "Starting the app..."
	go run ./cmd

## Очистка
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

## Docker: сборка
docker-build:
	@echo "Building Docker image..."
	docker build -t calculator-api .

## Docker: сборка без кэша
docker-build-nocache:
	@echo "Building Docker image without cache..."
	docker build --no-cache -t calculator-api .

## Docker: запуск
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 calculator-api

## Docker Compose: сборка
dc-build:
	@echo "Docker Compose build..."
	docker-compose build

## Docker Compose: сборка без кэша
dc-build-nocache:
	@echo "Docker Compose build without cache..."
	docker-compose build --no-cache

## Docker Compose: запуск
dc-up:
	@echo "Docker Compose up..."
	docker-compose up

## Docker Compose: остановка
dc-down:
	@echo "Docker Compose down..."
	docker-compose down

.PHONY: all build swagger test run clean \
        docker-build docker-build-nocache docker-run \
        dc-up dc-down dc-build dc-build-nocache