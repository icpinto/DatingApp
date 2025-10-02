#!/bin/bash
set -euo pipefail

MIGRATIONS_DIR=${MIGRATIONS_DIR:-/migrations}

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
        echo "Migrations directory $MIGRATIONS_DIR not found" >&2
        exit 1
fi

run_sql() {
        local file="$1"
        if [[ ! -f "$file" ]];
        then
                echo "Migration file $file not found" >&2
                exit 1
        fi

        echo "Applying migration: $(basename "$file")"
        psql -v ON_ERROR_STOP=1 \
                --username "$POSTGRES_USER" \
                --dbname "$POSTGRES_DB" \
                -f "$file"
}

shopt -s nullglob
readarray -t migrations < <(find "$MIGRATIONS_DIR" -maxdepth 1 -type f -name '*.sql' -print | sort)

if [[ ${#migrations[@]} -eq 0 ]]; then
        echo "No migration files found in $MIGRATIONS_DIR" >&2
        exit 1
fi

for migration in "${migrations[@]}"; do
        run_sql "$migration"
done

shopt -u nullglob

echo "All migrations applied successfully."
