# note : dockerfile run from the perspective of the docker-compose.yml file

# --- Build Stage (compile Go app) ---
FROM golang:1.22 AS builder

WORKDIR /app

# install go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (targeting correct package)
RUN go build -o nbaoverunders ./cmd/web
# at this point, we've compiled the Go app into a binary called nbaoverunders

# --- Final Stage (copy binary we built into smaller runtime image) ---
# base image
FROM debian:bookworm-slim

# Install MySQL client (needed for health checks)
RUN apt-get update && apt-get install -y default-mysql-client netcat-openbsd && rm -rf /var/lib/apt/lists/*

# set working directory
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/nbaoverunders .

# Copy the template files into the final image
COPY --from=builder /app/ui /app/ui

# Expose the port the app runs on
EXPOSE 4000

# Run the Go app
CMD ["./nbaoverunders"]

# docker build --platform linux/amd64 -t jbh14/nbaoverunders_amd:11 .
# docker push jbh14/nbaoverunders_amd:11