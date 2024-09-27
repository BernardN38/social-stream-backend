AUTH_BINARY=authApp



## up_build: stops docker-compose (if running), builds all projects, and starts docker compose
up_build: build_auth
	@echo "Stopping docker images (if running...)"
	docker compose down 
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d nginx postgres rabbitmq auth-service
	@echo "Docker images built and started!"
## up_build: stops docker-compose (if running), builds all projects, and starts docker compose
up_build_services: build_auth
	@echo "Stopping docker images (if running...)"
	docker compose down nginx auth-service 
	@echo "Building (when required) and starting docker images..."
	docker compose up --build nginx auth-service
	@echo "Docker images built and started!"
## build_auth: builds the authentication service binary as a linux executable
build_auth:
	@echo "Building auth service binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"
