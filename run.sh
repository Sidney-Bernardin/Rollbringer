function migrate-postgres {
    migrate \
        -path internal/repos/sql/Migrations \
        -database "postgres://$APP_POSTGRES_HOST:$APP_POSTGRES_PORT/rollbringer" \
        $@
}

function migrate-postgres-create {
    migrate-postgres create \
        -ext sql \
        -dir internal/repos/sql/Migrations \
        $@
}

###

function sqlc-generate {
    sqlc generate -f server/repositories/sql/sqlc.yaml
}

###

$@
