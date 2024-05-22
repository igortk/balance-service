FROM golang:1.20

WORKDIR /balance-service

COPY . .

RUN go build -o balance-service

EXPOSE 8080

CMD ["./balance-service"]