golang:
	air

parcel:
	npx parcel watch

db:
	docker run \
		-e POSTGRES_PASSWORD="password123" \
		postgres:latest 
		-name rollbringer-db \
		-p 5432:5432 \
		-v /tmp/pg-data:/var/lib/postgresql/data \
