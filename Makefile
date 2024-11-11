
.PHONY: docker-up
docker-up:
	docker compose -f ./docker-compose.local.yaml up --build --wait

.PHONY: docker-down
docker-down:
	docker compose -f ./docker-compose.local.yaml down > /dev/null 2>&1

.PHONY: docker-restart
docker-restart:
	make docker-down
	make docker-up

.PHONY: init-test
init-test:
	go install go.uber.org/mock/mockgen@latest
	go get go.uber.org/mock/mockgen/model

.PHONY: test-exercise
test-exercise:
	go clean -testcache && go test -race ./exercise/...

.PHONY: test
test:
	go test -race -short -coverpkg=training/app/... -coverprofile coverage.out ./app/...

.PHONY: database
database:
	go run cmd/database/main.go

.PHONY: api
api:
	go run cmd/api/main.go

.PHONY: producer
producer:
	go run cmd/producer/main.go

.PHONY: consumer
consumer:
	go run cmd/consumer/main.go

.PHONY: redis
redis:
	go run cmd/redis/main.go

.PHONY: gen-mock
gen-mock:
	find . -name "mock.go" -type f -delete
	mockgen -destination app/product/mock.go -package=product training/app/product Repository
	mockgen -destination app/user/mock.go -package=user training/app/user Repository
