# --- Stage 1: Build the React Frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ .
# Inject VITE_API_URL=/api for production
RUN VITE_API_URL=/api npm run build

# --- Stage 2: Build the Go Backend ---
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
# Copy go mod files from server directory
COPY server/go.mod server/go.sum ./server/
WORKDIR /app/server
RUN go mod download

# Copy the rest of the server source
WORKDIR /app
COPY server/ ./server/

# Copy the frontend build artifacts into the server directory for Go embedding
COPY --from=frontend-builder /app/web/dist ./server/webdist

# Build the Go application
WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- Stage 3: Final Runtime Image ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/server/server .
EXPOSE 8080
CMD ["./server"]
