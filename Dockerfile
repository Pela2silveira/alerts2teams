# Use the official Golang image to create a build artifact
FROM golang:1.19

# Set the working directory
WORKDIR /app

# Copy the local package files to the container's workspace
COPY *.go ./

# Build the Go application
RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -o /alert2teams


# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["/alert2teams"]
