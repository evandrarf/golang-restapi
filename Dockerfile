# Use the official Go image as the base image
FROM golang:1.21-bullseye

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o build

# Make the app executable
RUN chmod +x build

# Expose the port that the app listens on
EXPOSE 8080

# Set the command to run the app when the container starts
CMD ["./build"]
