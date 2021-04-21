FROM golang:1.14-alpine as builder
WORKDIR /src/go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .



RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client ./cmd
CMD ["./client","--port=10001","--server_addr=localhost:10000","--metric_port=4001"]




