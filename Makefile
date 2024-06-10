api-run:
	go run cmd/*
docker-build:
	docker build -t marketplace-api .
docker-run:
	docker run --name=marketplace-web-app -p 4000:4000 marketplace-api

docker-compose-build-up:
	docker-compose up --build marketplace-api

docker-compose-up:
	docker-compose up marketplace-api