golang:
	air

parcel:
	npx parcel watch --no-cache

db:
	docker run \
		--name rollbringer_db \
		-e POSTGRES_PASSWORD="password123" \
		-p 5432:5432 \
		-v ./tmp/pg-data:/var/lib/postgresql/data \
		postgres:latest
