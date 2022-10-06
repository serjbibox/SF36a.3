init: docker-down-clear \
	docker-pull docker-build docker-up \

docker-down-clear:
	docker-compose down -v --remove-orphans

docker-pull:
	docker-compose pull
	
docker-build:
	docker-compose build --pull
	
docker-up:
	docker-compose up -d