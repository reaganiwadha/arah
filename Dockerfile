FROM golang:1.18 AS builder

ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o arah

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV GIN_MODE=release
WORKDIR /root/
COPY templates/ ./templates/
COPY --from=builder /app/arah .
CMD ["./arah"]