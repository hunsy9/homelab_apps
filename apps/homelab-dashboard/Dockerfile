# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY backend/ ./
RUN go mod tidy
COPY --from=frontend-builder /app/backend/static ./static
RUN CGO_ENABLED=0 GOOS=linux go build -o /homelab-dashboard ./cmd/main.go

# Final image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /homelab-dashboard .
COPY --from=backend-builder /app/static ./static
EXPOSE 8080
CMD ["./homelab-dashboard"]
