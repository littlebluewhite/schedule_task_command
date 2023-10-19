swagger-docs:
	swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/api

run:
	go run cmd/api/main.go

gen:
	go run cmd/generate/main.go

api:
	go run cmd/api/main.go

up:
	go run cmd/migrate/main.go -up

up_test:
	go run cmd/migrate/main.go -up -test

down:
	go run cmd/migrate/main.go -down

down_test:
	go run cmd/migrate/main.go -down -test

to:
	go run cmd/migrate/main.go -to -version 1

swagger-docs-generate:
	swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/api

build-m:
	docker build -t schedule-migrate:latest -f deploy/migrate/linux/Dockerfile .

build-a:
	docker build -t schedule:latest -f deploy/api/linux/Dockerfile .

run-m:
	docker run --name schedule-migrate --rm --network="host" schedule-migrate:latest

run-a:
	docker run --name schedule -p 5487:5487 --network="host" -v ${PWD}/docker/log:/app/log schedule:latest


# docker exec -it <container-id> sh