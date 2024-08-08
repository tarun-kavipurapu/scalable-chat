# Define default port and server
PORT ?= 8080

# Targets
.PHONY: chat migrateup migratedown

# Run the chat server
chat:
	go run cmd/app/main.go -port=$(PORT)

# Run database migrations up
migrateup:
	go run cmd/migration/init.sql.go up

# Run database migrations down
migratedown:
	go run cmd/migration/init.sql.go down

# Run filestructure tool
filestructure:
	go run cmd/filestructure/main.go
