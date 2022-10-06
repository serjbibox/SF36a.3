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

stop:
	docker-compose down
run:
	docker-compose up -d	

test:
	docker-compose -f docker-compose_test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose_test.yml down --volumes

test-db-up:
	docker-compose -f docker-compose_test.yml up --build db

test-db-down:
	docker-compose -f docker-compose_test.yml down --volumes db