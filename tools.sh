function watch-server {
    air $@
}

function watch-web {
    npx parcel watch $@
}

#####

function sqlc {
    sqlc \
        -f server/repositories/sql/sqlc.yaml \
        $@
}

function migrate-postgres {
    migrate \
        -path server/repositories/sql/Migrations \
        -database $APP_POSTGRES_URL \
        $@
}

function migrate-postgres-create {
    migrate-postgres create \
        -ext sql \
        -dir server/repositories/sql/Migrations \
        $@
}

#####

function docker-compose {
    docker compose \
        -f Docker/docker-compose.yaml \
        --project-directory . \
        $@
}

$@
