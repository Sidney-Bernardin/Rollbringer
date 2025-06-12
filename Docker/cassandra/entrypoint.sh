#!/bin/bash

(
    # Wait for Cassandra.
    until (cqlsh -e 'describe cluster' > /dev/null 2>&1) do
        echo "[$0] Waiting to apply init queries..."
        sleep 1
    done

    # Apply the init queries to Cassandra.
    cqlsh -f /setup/init.cql
    echo "[$0] Done applying init queries."
) & docker-entrypoint.sh
