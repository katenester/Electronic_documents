build:
	docker-compose build doc-app

run:
	docker-compose up doc-app

test:
	go test -v ./...

migrate_up:
	migrate -path ./schema -database 'postgres://postgres:Katy314@0.0.0.0:5436/postgres?sslmode=disable' up

migrate_down:
	migrate -path ./schema -database 'postgres://postgres:Katy314@0.0.0.0:5436/postgres?sslmode=disable' down

swag:
	swag init -g cmd/main.go

stop:
	docker-compose stop