

# automatically builds and runs this in a docker container
deploy:
	docker-compose up -d --build

redeploy:
	docker-compose up -d --force-recreate --build


build:
	go build .