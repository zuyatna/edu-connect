FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o institution-service

EXPOSE 8081 50052

ENV GRPC_ENDPOINT=0.0.0.0
ENV GRPC_PORT=50052

CMD ["./institution-service"]