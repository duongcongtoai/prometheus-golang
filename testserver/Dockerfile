FROM golang:1.14-alpine as builder
WORKDIR /src/go

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd
CMD ["./server","--port=10000","--metric_port=4000"]




