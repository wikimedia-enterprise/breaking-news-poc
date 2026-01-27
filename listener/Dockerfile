FROM golang:1.23-alpine

# Switch to /app dir
WORKDIR /app

# Cache dependencies download
COPY go.mod go.sum ./
RUN go mod download

# Copy the project files
COPY . .

# Build the binary
RUN go build -o main listener/*.go

# Set the binary as the CMD of the container
CMD ["./main"]