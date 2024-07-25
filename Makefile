db_port := 8091

up:
	docker-compose up -d --build

down:
	docker-compose down

infra:
	docker-compose up -d db

migrate-new:
	goose -dir ./migrations create $(name) sql

migrate-up:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=dev host=localhost port=${db_port} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=dev host=localhost port=${db_port} sslmode=disable" down
