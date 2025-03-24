FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o transaction-service

EXPOSE 8082 50053

CMD ["./transaction-service"]