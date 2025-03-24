FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o user-service-example

EXPOSE 8080 50051

ENV GRPC_ENDPOINT=0.0.0.0
ENV GRPC_PORT=50051

CMD ["./user-service-example"]