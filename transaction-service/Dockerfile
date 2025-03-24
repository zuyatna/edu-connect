FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o transaction-service

EXPOSE 8082 50053

ENV GRPC_USER_ENDPOINT=user-service
ENV GRPC_USER_PORT=50051
ENV GRPC_INSTITUTION_ENDPOINT=institution-service
ENV GRPC_INSTITUTION_PORT=50052

CMD ["./transaction-service"]