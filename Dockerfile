# Multi-stage build for Ryanair Flight Planner MCP Server
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Enable direct proxy and disable CGO for pure static binary
ENV CGO_ENABLED=0
ENV GOPROXY=direct

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o ryanair-mcp-server .

# Production minimal runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/ryanair-mcp-server /app/ryanair-mcp-server

# Standard environment variables
ENV RYANAIR_CURRENCY=EUR
ENV RYANAIR_RADIUS_KM=150
ENV MCP_TRANSPORT=stdio

ENTRYPOINT ["/app/ryanair-mcp-server"]
