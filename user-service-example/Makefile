.PHONY: proto mockgen test

proto:
	protoc --go_out=. --go-grpc_out=. pb/*.proto

mockgen:
	mockgen -destination=./mocks/mock_user_repository.go -package=mocks github.com/zuyatna/edu-connect/user-service/repository IUserRepository \
	&& mockgen -destination=./mocks/mock_user_usecase.go -package=mocks github.com/zuyatna/edu-connect/user-service/usecase IUserUsecase \
	&& mockgen -destination=./mocks/mock_user_handler.go -package=mocks github.com/zuyatna/edu-connect/user-service/handler IUserHandler

test:
	go test -cover -v ./...