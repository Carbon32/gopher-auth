#!/bin/bash

case "$1" in
    create)
        if [ -z "$2" ]; then
            echo "Usage: ./migrate.sh create <migration_name>"
            exit 1
        fi

        migrate create -ext sql -dir cmd/migrate/migrations -seq "$2"
        ;;
    up)
        go run ./cmd/migrate up
        ;;

    down)
        go run ./cmd/migrate down
        ;;

    *)
        echo "Usage: ./migrate.sh (create/up/down)"
        exit 1
        ;;
esac