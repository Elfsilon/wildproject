FROM golang:1.22.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./app ./cmd/api-server/main.go

EXPOSE 8000

CMD ["./app"]