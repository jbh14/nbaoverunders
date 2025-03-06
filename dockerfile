# note : dockerfile run from the perspective of the docker-compose.yml file

# --- Build Stage ---
FROM golang:1.22 AS builder

WORKDIR /app

# install go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (targeting correct package)
RUN go build -o nbaoverunders ./cmd/web

# --- Final Stage ---
# base image
FROM debian:bookworm-slim

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