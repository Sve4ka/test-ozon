memory:
	docker-compose --env-file .env.memory up --build

postgres:
	docker-compose --env-file .env.postgres up --build

cache:
	docker-compose --env-file .env.redis.postgres up --build

swag:
	swag init --parseDependency -d ./internal/delivery/http/handler -g ../../../../cmd/main.go -o ./internal/docs

mock:
	mockery --name=Storage --dir=./internal/storage --output=./internal/mocks --outpkg=mocks --filename=storage.go
	mockery --name=Cache --dir=./internal/cache --output=./internal/mocks --outpkg=mocks --filename=cache.go
	mockery --name=Generator --dir=./internal/link --output=./internal/mocks --outpkg=mocks --filename=generator.go
