FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go
RUN ls -la  # Verifica se o executável foi criado

FROM alpine:3.20.1 AS prod

WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=build /app/main /app/main
COPY --from=build /app/migrations /app/migrations/
COPY --from=build /app/scripts /app/scripts/
COPY --from=build /app/.env.example /app/.env

RUN chmod +x /app/main
RUN chmod +x /app/scripts/init-db.sh

EXPOSE ${PORT}
CMD ["./main"]


