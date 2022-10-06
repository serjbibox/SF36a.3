init: docker-down-clear \
	docker-pull docker-build docker-up \

docker-down-clear:
	docker-compose -f docker-compose.yml down -v --remove-orphans

docker-pull:
	docker-compose -f docker-compose.yml pull
	
docker-build:
	docker-compose -f docker-compose.yml build --pull

docker-up:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml down
run:
	docker-compose -f docker-compose.yml up -d	

test:
	docker-compose -f docker-compose_test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose_test.yml down --volumes

test-db-up:
	docker-compose -f docker-compose_test.yml up --build db

test-db-down:
	docker-compose -f docker-compose_test.yml down --volumes db