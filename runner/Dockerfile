# Stage 1: Build the Go application
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o roadmap .

# Stage 2: Create a smaller image for runtime
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the binary from the builder stage
COPY --from=builder /usr/src/app/roadmap .

# Copy the input file into the container
COPY input.txt /usr/src/app/input.txt
# Command to run the executable
CMD ["./roadmap", "input.txt"]
