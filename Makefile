

# automatically builds and runs this in a docker container
deploy:
	docker-compose up -d --build


build:
	go build .