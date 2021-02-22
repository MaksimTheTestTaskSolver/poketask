.PHONY: run build

build:
	docker build . -t poketask

run:
	docker run -p 8080:8080 poketask

start: build run