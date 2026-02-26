APP_NAME=alloy-ui
ENV=.env

.PHONY: build run stop clean restart logs

build:
	docker build --target dev -t $(APP_NAME) .

run:
	docker run -d --name $(APP_NAME) --env-file $(ENV) -p 3000:3000 $(APP_NAME)


stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true


restart: stop build run

clean: stop
	docker rmi $(APP_NAME) || true

logs:
	docker logs -f $(APP_NAME)
