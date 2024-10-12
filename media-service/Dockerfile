# Use a minimal Alpine base image
FROM alpine:latest

# Install necessary packages (if needed)
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the local target directory
COPY target/x86_64-unknown-linux-musl/release/media-service .

# Set the command to run your application
CMD ["./media-service"]
