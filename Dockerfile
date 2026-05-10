FROM node:20-alpine AS frontend
WORKDIR /app
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

FROM golang:alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=backend /server /app/server
COPY --from=frontend /app/dist /app/web/dist
ENV WEB_DIST=/app/web/dist
ENV DATA_DIR=/app/data
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
