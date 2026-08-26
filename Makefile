APP := storyshard-app
IMAGE := storyshard-api:latest
COMPOSE := docker compose -f deploy/docker-compose.prod.yml

.PHONY: build run test clean docker-build docker-prod docker-prod-api docker-migrate docker-prod-down docker-logs

build: clean
	go build -o $(APP) ./cmd/$(APP)

run: build
	./$(APP)

test:
	go test -v ./...

clean:
	rm -f $(APP)

docker-build:
	docker build -f deploy/Dockerfile -t $(IMAGE) .

docker-migrate:
	$(COMPOSE) --profile tools run --rm migrate

docker-prod:
	$(COMPOSE) up -d api web

docker-prod-api:
	$(COMPOSE) up -d api

docker-prod-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f api web
