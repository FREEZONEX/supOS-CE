# Use official Go image as build environment
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary system dependencies
RUN apk add --no-cache ca-certificates 

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application - static compilation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o backend ./backend.go

# Use lightweight alpine image as runtime environment
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# Set working directory
WORKDIR /app

# Copy binary file from build stage
COPY --from=builder /app/backend .
RUN chmod +x ./backend

# Copy configuration files from build stage
COPY --from=builder /app/etc ./etc

# Expose port
EXPOSE 8080

# Start application
CMD ["./backend"]
