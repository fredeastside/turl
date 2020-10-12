build: ; docker-compose up -d --build

start: build

down: ; docker-compose down

up: ; docker-compose up -d

stop: ; docker-compose stop

test: ; docker-compose -p turl_tests -f docker-compose.test.yml up --build --abort-on-container-exit