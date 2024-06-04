build:
	cd ./server
	docker-compose down
	@echo "Stopping and removing existing containers..."
	docker-compose up --build -d
	@echo "Successfully built and started new containers!"
	@echo "Client available at http://localhost:3000"
	@echo "Server API available at http://localhost:8080"

up:
	cd ./server
	@echo "Starting containers in detached mode..."
	docker-compose up --build -d
	@echo "Containers started successfully!"
	@echo "Client available at http://localhost:3000"
	@echo "Server API available at http://localhost:8080"

down:
	docker-compose down
	@echo "Containers stopped and removed successfully!"

stop:
	docker-compose stop
	@echo "Containers stopped successfully!"

.PHONY: format
format:
	@echo "Formatted all Go Source Files..."
	cd ./server && gofmt -w .

.PHONY: test
test:
	@echo "Performing Unit Tests..."
	cd ./server && go test ./...
