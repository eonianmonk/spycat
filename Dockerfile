# Use the official Go 1.22 image as the base image
FROM golang:1.22-bullseye

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Set the GOARCH environment variable to amd64
ENV GOARCH=amd64

# Build the Go application
RUN go build -o spycat

RUN chmod +x spycat
# Set the entrypoint to run the binary
ENTRYPOINT ["./spycat"]