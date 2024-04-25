build:
	@docker build -t app:local .

up:
	@make build
	@docker-compose -f docker-compose.yml up -d

stop:
	@docker-compose -f docker-compose.yml stop

down:
	@docker-compose -f docker-compose.yml down