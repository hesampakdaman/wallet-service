FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main/main.go
RUN go build -o migrate ./cmd/migrate/migrate.go

EXPOSE 8080

CMD ["./main"]
