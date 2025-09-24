FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bookshop ./cmd/bookshop


FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /root/
COPY --from=builder /app/bookshop /usr/local/bin/
COPY config.yaml ./

ENV CONFIG_PATH=config.yaml
ENTRYPOINT ["bookshop"]
CMD ["--help"]