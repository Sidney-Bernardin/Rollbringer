function watch-server {
    air $@
}

function watch-web {
    rm -r .parcel-cache
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
    migrate create \
        -ext sql \
        -dir server/repositories/sql/Migrations \
        $@
}

function migrate-cassandra {
    migrate \
        -path server/repositories/cql/Migrations \
        -database "cassandra://$(cut -d ',' -f 1 <<< $APP_CASSANDRA_HOSTS)/rollbringer" \
        $@
}

function migrate-cassandra-create {
    migrate create \
        -ext sql \
        -dir server/repositories/cql/Migrations \
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
