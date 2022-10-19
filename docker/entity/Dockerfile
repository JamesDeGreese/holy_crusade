FROM golang:alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o entity_service_app ./cmd/entity/app
RUN go build -o entity_service_cli ./cmd/entity/cli
CMD ["./entity_service_app"]