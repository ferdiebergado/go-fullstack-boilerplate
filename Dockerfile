# Use the official Go image as the base image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install Air for hot reloading
RUN go install github.com/air-verse/air@v1.52.3

# Set environment variables for development
ENV GO_ENV=development

# Expose the port the app runs on
EXPOSE 8888

# Command to run the application with Air for hot reloading
CMD ["air"]
