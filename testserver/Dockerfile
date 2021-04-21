FROM golang:1.14-alpine as builder
WORKDIR /src/go
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd
CMD ["./server"]




