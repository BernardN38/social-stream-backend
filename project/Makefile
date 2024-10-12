AUTH_BINARY=authApp
USER_BINARY=userApp
MEDIA_BINARY=mediaApp
MEDIA_COMPRESSION_BINARY=mediaCompressionApp

## up_build: stops docker-compose (if running), builds all projects, and starts docker compose
up_build: build_auth build_media build_user build_media_compression
	@echo "Stopping docker images (if running...)"
	docker compose down 
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## up_build: stops docker-compose (if running), builds all projects, and starts docker compose
up_build_services: build_auth build_user build_media build_media_compression
	@echo "Stopping docker images (if running...)"
	docker compose down nginx auth-service user-service media-service media-compression-service
	@echo "Building (when required) and starting docker images..."
	docker compose up --build nginx auth-service user-service media-service media-compression-service
	@echo "Docker images built and started!"
	
## build_auth: builds the authentication service binary as a linux executable
build_auth:
	@echo "Building auth service binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"


## build_user: builds the user service binary as a linux executable
build_user:
	@echo "Building user service binary..."
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"


## build_media: builds the media service binary as a linux executable
build_media:
	@echo "Building media service binary..."
	cd ../media-service && cargo build --release --target=x86_64-unknown-linux-musl
	@echo "Done!"


## build_media: builds the media service binary as a linux executable
build_media_compression:
	@echo "Building media compression service binary..."
	cd ../media-compression-service && cargo build --release --target=x86_64-unknown-linux-musl
	@echo "Done!"


