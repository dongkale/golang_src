docker:
	docker build -t my-service:latest .

build:
	go build -o my-service cmd/main.go

run:
	go run cmd/main.go

migration-status:
	migrate -database "mysql://root:root@tcp(localhost:3306)/my_boilerplate_service?multiStatements=true" -path migrations version

migrate:
	migrate -database "mysql://root:root@tcp(localhost:3306)/my_boilerplate_service?multiStatements=true" -path migrations up

unmigrate:
	migrate -database "mysql://root:root@tcp(localhost:3306)/my_boilerplate_service?multiStatements=true" -path migrations down

create-migration:
	migrate create -ext sql -dir migrations -seq $(name)