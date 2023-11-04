# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .


# Build the Go application
RUN go build -o main .

# Expose the port on which the application will run
EXPOSE 8080

# Run the application
CMD ["./main"]
