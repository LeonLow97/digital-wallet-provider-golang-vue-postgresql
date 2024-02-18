build:
	cd ./server
	docker-compose down
	@echo "Building and starting docker images"
	docker-compose up --build -d
	@echo "Docker images built and started!"
	@echo "Client running on port 3000"
	@echo "Server running on port 8080"

down:
	docker-compose down