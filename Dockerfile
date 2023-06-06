FROM golang:latest

LABEL maintainer="Thomas Chardonnens <thomas.chardonnens@berkeley.edu>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 8080

CMD ["./main"]

ENV DB_PATH=/home/debian/t3/t3.db
