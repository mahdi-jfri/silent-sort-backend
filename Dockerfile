
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/app .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/app /

EXPOSE 8080

ENTRYPOINT ["/app"]