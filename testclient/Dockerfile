FROM golang:1.14-alpine as builder
WORKDIR /src/go
COPY . .

RUN go get -u github.com/tsenart/vegeta

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client ./cmd
CMD ["./client"]




