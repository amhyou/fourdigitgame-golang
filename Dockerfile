# Start with a base Go image to build the binary
FROM golang:1.22-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o /app/server .

# Start a new stage from scratch
FROM alpine:latest

# Copy the pre-built binary file from the previous stage
COPY --from=build /app/server /app/server

# Expose port 5000 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["/app/server"]
