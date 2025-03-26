.PHONY: proto mockgen test

proto:
	protoc --go_out=. --go-grpc_out=. pb/*.proto

mockgen:
	mockgen -destination=./mocks/mock_institution_repository.go -package=mocks institution-service/repository IInstitutionRepository \
	&& mockgen -destination=./mocks/mock_post_repository.go -package=mocks institution-service/repository IPostRepository \
	&& mockgen -destination=./mocks/mock_fund_collect_repository.go -package=mocks institution-service/repository IFundCollectRepository \
	&& mockgen -destination=./mocks/mock_institution_usecase.go -package=mocks institution-service/usecase IInstitutionUsecase \
	&& mockgen -destination=./mocks/mock_post_usecase.go -package=mocks institution-service/usecase IPostUsecase

test:
	go test -cover -v ./...