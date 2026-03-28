IMAGE_NAME := todod
DEV_PORT   := 8085

.PHONY: build run

build:
	docker build --target release-todod -t $(IMAGE_NAME) .

run: build
	docker run --rm -p $(DEV_PORT):8080 $(IMAGE_NAME)
