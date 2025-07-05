
build: 
	go build -o app ./cmd/api

run: build 
	./app 

clean: 
	rm -f app

migrate-up:
	go run cmd/migrate/migration.go -direction=up

migrate-down:
	go run cmd/migrate/migration.go -direction=down
	
.PHONY: docs
docs:
	swag init --dir "./cmd/api,./internal/data" 

