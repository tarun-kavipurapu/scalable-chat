
chat:
	go run /cmd/app/main.go


migrateup:
	go run /cmd/migration/init.sql.go up

migratedown:
	go run /cmd/migration/init.sql.go down

filestructure:
	go run /cmd/filestructure/main.go 