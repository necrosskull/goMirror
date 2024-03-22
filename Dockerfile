FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -tags=jsoniter -o main main.go

FROM alpine:latest
WORKDIR /bin
COPY --from=builder /app/main .
CMD [ "/bin/main" ]