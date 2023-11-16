FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o main main.go

EXPOSE 8085

CMD ["/app/main"]