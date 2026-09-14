FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o proxy ./cmd/proxy

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/proxy /proxy
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/proxy"]
