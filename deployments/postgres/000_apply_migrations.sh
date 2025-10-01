#!/bin/bash
set -euo pipefail

MIGRATIONS_DIR=${MIGRATIONS_DIR:-/migrations}

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

run_sql "$MIGRATIONS_DIR/init_schema.sql"
run_sql "$MIGRATIONS_DIR/add_username_columns.sql"
run_sql "$MIGRATIONS_DIR/add_description_to_friend_requests.sql"
run_sql "$MIGRATIONS_DIR/add_profile_image_column.sql"
run_sql "$MIGRATIONS_DIR/add_profile_verification_columns.sql"
run_sql "$MIGRATIONS_DIR/add_profile_sync_outbox_table.sql"
run_sql "$MIGRATIONS_DIR/update_profiles_schema.sql"
run_sql "$MIGRATIONS_DIR/update_conversation_id_uuid.sql"

echo "All migrations applied successfully."
