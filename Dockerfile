FROM golang:1.20-alpine as builder
WORKDIR /app/
COPY . .
RUN apk add --no-cache --update gcc g++ && CGO_ENABLED=1 go build -o bin/paste main.go

FROM alpine:3.17
WORKDIR /app/
COPY --from=builder /app/bin/paste .
COPY web/templates ./web/templates
ENTRYPOINT ["./paste"]
