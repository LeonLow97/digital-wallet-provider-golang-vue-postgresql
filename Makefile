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