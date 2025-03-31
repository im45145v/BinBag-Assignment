# Use golang:1.23-alpine as the base image
FROM golang:1.23-alpine

# Set the working directory to /app
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Run go mod download to install dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Run go mod tidy to clean up and synchronize dependencies
RUN go mod tidy

# Run go build -o main . to build the application
RUN go build -o main .

# Set the entry point to ./main
ENTRYPOINT ["./main"]
