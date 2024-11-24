FROM golang:1.23-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o mintestApp .

RUN chmod +x /app/mintestApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mintestApp /app

CMD ["/app/mintestApp"]




