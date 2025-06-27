APP = personae-fasti
DOCKER_VOLUME=personae
DOCKER_NETWORK=personae
DOCKER_HOST=172.41.2.1
DOCKER_PORT=4121  # The port your app listens on inside the container
HOST_PORT=4121  # The port you want to expose on the host

DOCKER_FILE_SERVER_NETWORK=file-server
DOCKER_FILE_SERVER_HOST=172.72.2.1


build: clean
	@go build -o bin/$(APP)

run: build
	@./bin/$(APP)

docker:
	@docker image build -t $(APP):v$(cv) .
	@docker container create --name $(APP) -v $(DOCKER_VOLUME):/app/mnt $(APP):v$(cv)
	@make docker-run	

docker-run:
	@docker container start $(APP)
	@docker network connect --ip $(DOCKER_HOST) $(DOCKER_NETWORK) $(APP)
	@docker network connect --ip $(DOCKER_FILE_SERVER_HOST) $(DOCKER_FILE_SERVER_NETWORK) $(APP)

docker-stop:
	@docker container stop $(APP)

docker-rm: 
	@docker container stop $(APP)
	@docker container rm $(APP)
	@docker image rm $(APP):v$(cv)

test: clean
	@go test -v ./...

clean:
	@rm -rf bin/*