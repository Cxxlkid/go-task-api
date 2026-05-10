# Stage 1 — Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# On copie les fichiers de dépendances en premier
# Docker met en cache cette couche si go.mod/go.sum n'ont pas changé
COPY go.mod go.sum ./
RUN go mod download

# On copie le reste du code
COPY . .

# On compile l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

# Stage 2 — Run
FROM alpine:latest

WORKDIR /app

# On copie uniquement le binaire compilé du stage précédent
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]