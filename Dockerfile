FROM golang:1.20.3-alpine3.17

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /t3-api

EXPOSE 8080

CMD ["/t3-api"]