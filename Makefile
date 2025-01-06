run:
	- go run .

migrations:
	- goose create $m sql -dir ./db/migrations -env ./.env

migrate_up:
	- goose up -dir ./db/migrations -env ./.env

migrate_down:
	- goose down -dir ./db/migrations -env ./.env

db_status:
	- goose status -env ./.env -dir ./db/migrations

air:
	- air

swag:
	- swag init -g api/api.go