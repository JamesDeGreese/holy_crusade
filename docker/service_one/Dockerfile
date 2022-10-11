FROM golang:alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o service_one ./cmd/service_one
CMD ["./service_one"]