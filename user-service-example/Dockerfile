FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o user-service

EXPOSE 8080

CMD ["./user-service"]