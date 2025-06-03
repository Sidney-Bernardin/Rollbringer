function docker-compose {
    docker compose \
        -f Docker/docker-compose.yaml \
        --project-directory . \
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

function sqlc-generate {
    sqlc generate -f server/repositories/sql/sqlc.yaml
}

###

$@
