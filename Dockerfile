# ---------- build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Dependencias del sistema
RUN apk add --no-cache ca-certificates mysql-client

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o daycare-api ./cmd/api

# ---------- runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Certificados SSL (importante para JWT, HTTPS, etc.)
RUN apk add --no-cache ca-certificates mysql-client

# Copiar binario compilado
COPY --from=builder /app/daycare-api /app/daycare-api

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Puerto donde escucha la API
EXPOSE 8080

CMD ["/entrypoint.sh"]

# Ejecutar la API
#CMD ["./daycare-api"]

