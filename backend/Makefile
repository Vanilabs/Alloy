APP_NAME=alloy-api
PORT=8080
ENV=.env

.PHONY: build run stop clean test restart logs

test:
	go test ./tests/... -v

build: test
	docker build -t $(APP_NAME) .

run:
	docker run -d --name $(APP_NAME) --env-file $(ENV) -p $(PORT):8080 $(APP_NAME)


stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true


restart: stop build run

clean: stop
	docker rmi $(APP_NAME) || true

logs:
	docker logs -f $(APP_NAME)

run-postgres-only:
	docker run -d \
	  --name alloy-postgres \
	  --env-file $(ENV) \
	  -p 5432:5432 \
	  postgres:16.11-alpine3.23
