FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o institution-service

EXPOSE 8080
EXPOSE 50051

CMD ["./institution-service"]