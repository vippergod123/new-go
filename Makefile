postgres:
	docker run --network bank-network --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres:15.4 
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate
server:
	go run main.go
mock:
	mockgen -destination db/.mock/store.go -package mockdb github.com/vipeergod123/simple_bank/db/sqlc Store
test:
	go test -v -cover ./...
.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1