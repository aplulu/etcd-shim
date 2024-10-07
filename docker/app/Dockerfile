FROM golang:1.23.2-alpine as builder

WORKDIR /go/src/app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/serve/main.go

FROM alpine:latest

COPY --from=builder /go/bin/app /

CMD ["/app"]
